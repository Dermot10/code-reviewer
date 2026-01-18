import json
import signal
import os
import redis.asyncio as redis
import asyncio
from pydantic import TypeAdapter
from backend_python.schemas.dto.task import Task
from backend_python.worker.processor import process_review_task, process_enhance_task
from backend_python.logger import get_logger
from dotenv import load_dotenv
load_dotenv()

logger = get_logger("Ai-Worker")

r = redis.Redis(host=os.getenv("REDIS_HOST"), port=os.getenv("REDIS_PORT"), db=0, decode_responses=True)
QUEUE_KEY = "queue:tasks"
task_adapter = TypeAdapter(Task)

stop_event = asyncio.Event()


# TODO - 
# deprecate handler logic, current purpose - sync testing
# move export functionality to Go service


async def handle_task(task_dict): 
    try: 
        task = task_adapter.validate_python(task_dict)
        lock_key = f"{task.type}:{getattr(task, f'{task.type}_id')}:lock"

        # if not exists set 
        if not await r.set(lock_key, "1", nx=True, ex=os.getenv("LOCK_EXPIRY")):
            logger.info("Task already in progress, skipping", extra={"task": task_dict})
            return
        
        logger.info("Processing Task", extra={"task": task.type, "task_id": getattr(task, f"{task.type}_id")})

        if task.type == "review": 
            # ai agents
            result = await process_review_task(task)
            result_key = f"review:{task.review_id}:result"

            # store in result in redis 
            await r.set(result_key, json.dumps(result))
            await r.publish("review.completed", json.dumps({"review_id": task.review_id}))

        elif task.type == "enhance":

            result = await process_enhance_task(task)
            result_key = f"enhancement:{task.enhancement_id}:result"

            await r.set(result_key, json.dumps(result))
            await r.publish("enhancement.completed", json.dumps({"enhancement_id": task.enhancement_id}))
        else:  
            logger.warning("Unknown task type", etxra={"task": task_dict})
            return 

        logger.info("Task completed", extra={"task":task_dict})

        await r.delete(lock_key)

    except Exception as e:
        logger.exception("Failed to process task", exc_info=e)



async def worker_loop(): 
    logger.info("Worker started, waiting for tasks...")

    while not stop_event.is_set():
        try: 
            # block if empty, pop if available
            item = await r.brpop(os.getenv("QUEUE_KEY"), timeout=5)
            if item: 
                _, task_data = item
                asyncio.create_task(handle_task(json.loads(task_data)))
        except Exception as e:
            logger.exception("Error reading from Redis", exc_info=e)
            await asyncio.sleep(1)


async def main(): 
    loop = asyncio.get_running_loop()
    loop.add_signal_handler(signal.SIGINT, stop_event.set)
    loop.add_signal_handler(signal.SIGTERM, stop_event.set)

    await worker_loop()
            
if __name__=="__main__": 
    asyncio.run(main())
