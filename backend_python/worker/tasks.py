import json
import asyncio
from backend_python.schemas.dto.task import Task
from backend_python.worker.processor import process_review_task, process_enhance_task, process_assistant_task



BATCH_SIZE = 120
FLUSH_INTERVAL = 0.1 # 100 ms

async def handle_review_task(r, task: Task): 
    result = await process_review_task(task)

    result_key = f"review:{task.review_id}:result"

    await r.set(result_key, json.dumps(result))

    await r.publish("review.completed", json.dumps({"review_id": task.review_id}))


async def handle_enhance_task(r, task: Task): 
    result = await process_enhance_task(task)

    result_key = f"enhancement:{task.enhancement_id}:result"

    await r.set(result_key, json.dumps(result))

    await r.publish("enhancement.completed", json.dumps({"enhancement_id": task.enhancement_id}))


async def handle_assistant_task(r, task: Task): 

    result = ""
    buffer = []
    last_flush = asyncio.get_event_loop().time()
    
    result_key = f"assistant:{task.conversation_id}:result"


    async for chunk in process_assistant_task(task):
        result += chunk
        buffer.append(chunk)
       
        #check for deprecation, may require updating           
        now = asyncio.get_event_loop().time()
        
        # multi flush conditions for to ensure streaming -> UI stays performant
        should_flush = (
            len(buffer) >= BATCH_SIZE or 
            (now - last_flush) >= FLUSH_INTERVAL 
        )

        if not should_flush: 
            continue
        
        payload = "".join(buffer)

        await r.publish(
            "assistant.events",
            json.dumps({
                "type": "assistant.stream", 
                "payload": {
                    "conversation_id": task.conversation_id, 
                    "chunk": payload,
                    "done": False
                }
            })
        )

        buffer.clear()
        last_flush = now
    
    # flush remaining buffer 
    if buffer: 
        payload = "".join(buffer)

        await r.publish(
            "assistant.events",
            json.dumps({
                "type": "assistant.stream",
                "payload": {
                    "conversation_id": task.conversation_id,
                    "chunk": payload,
                    "done": False
                }
            })
        )

    await r.set(result_key, result)

    await r.publish("assistant.events", json.dumps({
        "type": "assistant.stream",
        "payload": {
            "conversation_id": task.conversation_id,
            "chunk": result, 
            "done": True
            }
        })
    )