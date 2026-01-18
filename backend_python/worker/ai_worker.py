import json
import redis
import asyncio
import sys
from pathlib import Path
from pydantic import TypeAdapter
from backend_python.processing.postprocessing import postprocess_review, postprocess_enhanced
from backend_python.processing.preprocessing import extract_chunks
from backend_python.service.review_service import Execute_review
from backend_python.service.code_quality_service import Execute_enhance

project_root = Path(__file__).resolve().parent.parent  
sys.path.append(str(project_root))

from backend_python.schemas.dto.task import Task



r = redis.Redis(host="localhost", port=6379, db=0, decode_responses=True)
QUEUE_KEY = "queue:tasks"
task_adapter = TypeAdapter(Task)


# TODO - 
# Replace with better logging 
# remove handler logic after moving export function to Go service
# test for e2e bugs with the ai service
# refactor and remove any bugs in processing logic

async def process_task(task: Task): 
    if task.type == "review":
        print(f"[REVIEW] Processing task: user={task.user_id}, review={task.review_id}, action={task.action}")
        result = f"AI result for review {task.review_id}"
        print(f"Result: {result}")

        chunked_context = extract_chunks(code=task.code)
        print("")
        print(chunked_context)
        print("")
        response = await Execute_review(chunked_context)
        return postprocess_review(response)


    elif task.type == "enhance":
        print(f"[ENHANCE] Processing task: user={task.user_id}, enhancement={task.enhancement_id}, action={task.action}")
        result = f"Enhanced code result for enhancement {task.enhancement_id}"
        print(f"Result: {result}")

        chunked_context = extract_chunks(code=task.code)
        print("")
        print(f"{chunked_context}+\n")
        print("")
        response = await Execute_enhance(chunked_context)
        print("----- The response from the code quality -----")
        print(f"{response}+\n")
        print("")
        return postprocess_enhanced(response)
    else:
        print(f"Unknown task type: {task}")

    

async def main():
    
    print("Worker started, waiting for tasks...")
    while True: 
        _, data = r.brpop(QUEUE_KEY)
        task_dict = json.loads(data)
        try:
            task = task_adapter.validate_python(task_dict)\
            
            lock_key = f"review:{task.review_id}:lock"
            if not r.set(lock_key, "1", nx=True, ex=600): 
                print(f"Review {task.review_id} already processed or in progress")
                continue
        
            result = await process_task(task)

            if task.type == "review":
                r.set(f"review:{task.review_id}:result", json.dumps(result))
                r.publish("review.completed", json.dumps({"review_id": task.review_id}))
                print(f"✓ Review {task.review_id} completed")
                
            elif task.type == "enhance":
                r.set(f"enhancement:{task.enhancement_id}:result", json.dumps(result))
                r.publish("enhancement.completed", json.dumps({"enhancement_id": task.enhancement_id}))
                print(f"✓ Enhancement {task.enhancement_id} completed")


        except Exception as e:
            print(f"failed to parse task: {e}")
            continue

if __name__=="__main__": 
    asyncio.run(main())
