import json
import redis
import sys
from pathlib import Path
from pydantic import parse_obj_as

# Add the project root (one level above 'backend_python') to sys.path
project_root = Path(__file__).resolve().parent.parent  
sys.path.append(str(project_root))

from backend_python.schemas.dto.task import Task



r = redis.Redis(host="localhost", port=6379, db=0, decode_responses=True)
QUEUE_KEY = "queue:tasks"

def process_task(task: Task): 
    if task.type == "review":
        print(f"[REVIEW] Processing task: user={task.user_id}, review={task.review_id}, action={task.action}")
        result = f"AI result for review {task.review_id}"
        print(f"Result: {result}")
    elif task.type == "enhance":
        print(f"[ENHANCE] Processing task: user={task.user_id}, enhancement={task.enhancement_id}, action={task.action}")
        result = f"Enhanced code result for enhancement {task.enhancement_id}"
        print(f"Result: {result}")
    else:
        print(f"Unknown task type: {task}")

def main():
    print("Worker started, waiting for tasks...")
    while True: 
        _, data = r.brpop(QUEUE_KEY)
        task_dict = json.loads(data)
        try:
            task = parse_obj_as(Task, task_dict)
        except Exception as e:
            print(f"failed to parse task: {e}")
            continue
        
        process_task(task)



if __name__=="__main__": 
    main()
