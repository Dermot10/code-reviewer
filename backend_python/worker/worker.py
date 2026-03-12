import json
import signal
import os
import redis.asyncio as redis
import asyncio
from pydantic import TypeAdapter
from backend_python.schemas.dto.task import Task
from backend_python.worker.tasks import handle_assistant_task, handle_enhance_task, handle_review_task
from backend_python.logger import get_logger
from dotenv import load_dotenv

load_dotenv()

logger = get_logger("Worker")

r = redis.Redis(
    host=os.getenv("REDIS_HOST", "localhost"),
    port=int(os.getenv("REDIS_PORT", 6379)),
    db=0,
    decode_responses=True
)


QUEUE_KEY = os.getenv("QUEUE_KEY", "queue:tasks")
LOCK_EXPIRY = int(os.getenv("LOCK_EXPIRY", 30))


SEMAPHORE = asyncio.Semaphore(10)
task_adapter = TypeAdapter(Task)
stop_event = asyncio.Event()


# TODO - 
# move export functionality to Go service

async def run_task(task_data): 
    """Semaphore wrapper to limit concurrent execution"""
    async with SEMAPHORE: 
        await handle_task(task_data)


async def handle_task(task_dict):

    task = None 
    lock_key = None
    lock_acquired = False

    try:
        task = task_adapter.validate_python(task_dict)

        # map task types to their unique ID fields
        id_field_map = {
            "review": "review_id",
            "enhance": "enhancement_id",
            "assistant": "conversation_id"
        }

        task_id_field = id_field_map.get(task.type)

        if not task_id_field:
            logger.warning("Unknown task type for lock", extra={"task": task_dict})
            return

        task_id_value = getattr(task, task_id_field, None)

        if task_id_value is None:
            logger.warning("Task missing ID field", extra={"task": task_dict})
            return

        lock_key = f"{task.type}:{task_id_value}:lock"
        
        lock_acquired = await r.set(
            lock_key, "1", nx=True, ex=int(os.getenv("LOCK_EXPIRY", 30))
        )

        # idempotency lock
        if not await r.set(lock_key, "1", nx=True, ex=int(os.getenv("LOCK_EXPIRY", 30))):
            logger.info("Task already in progress, skipping", extra={"task": task_dict})
            return

        logger.info(
            "Processing Task", 
            extra={"task": task.type, "task_id": task_id_value}
        )

        if task.type == "review":
            await handle_review_task(r, task)

        elif task.type == "enhance":
            await handle_enhance_task(r, task)

        elif task.type == "assistant":
            await handle_assistant_task(r, task)

        logger.info("Task completed", extra={"task": task_dict})
        

    except Exception as e:
        logger.exception("Failed to process task", exc_info=e)

        if task and task.type == "assistant": 
            await r.publish(
                "assistant.events", json.dumps({
                "type": "assistant.failed", 
                "payload": {
                    "conversation_id": getattr(task, "conversation_id", None), 
                    "error": str(e)
                }
            })
        )

    finally: 
        if lock_acquired: 
            await r.delete(lock_key)

async def worker_loop(): 
    logger.info("Worker started", extra={"queue": QUEUE_KEY})

    while not stop_event.is_set():
        try: 
            # block if empty, pop if available
            item = await r.brpop(QUEUE_KEY, timeout=5)
            if item: 
                _, task_data = item

                # deserialise / decode json to python dict (str)
                task_dict = json.loads(task_data)

                # spawn controlled task 
                asyncio.create_task(run_task(task_dict))
                
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
