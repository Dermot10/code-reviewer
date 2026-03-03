from backend_python.processing.postprocessing import postprocess_review, postprocess_enhanced
from backend_python.processing.preprocessing import extract_chunks
from backend_python.service.review_service import Execute_review
from backend_python.service.code_quality_service import Execute_enhance
from backend_python.schemas.dto.task import Task
import json



async def process_review_task(task: Task): 
    """Process a 'review' type task"""

    chunked_context = extract_chunks(code=task.code)
    raw_response = await Execute_review(chunked_context)
    return postprocess_review(raw_response)



async def process_enhance_task(task: Task): 
    """Process an 'enhancement' type task"""

    chunked_context = extract_chunks(code=task.code)
    raw_response = await Execute_enhance(chunked_context)
    return postprocess_enhanced(raw_response)



async def process_assistant_task(task:Task): 
    """Process a 'chat' type task"""
    full_response = ""

    async for chunk in stream_ai_model(task.prompt): 
        full_response += chunk
        

        #publish individual chunk
        await r.publish("assistant.events", json.dumps({
            "type": "assistant.chunk", 
            "user_id": task.user_id, 
            "conversation_id": task.conversation_id, 
            "chunk": chunk
        }))

    return full_response
