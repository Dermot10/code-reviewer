from backend_python.processing.postprocessing import postprocess_review, postprocess_enhanced
from backend_python.processing.preprocessing import extract_chunks
from backend_python.service.review_service import Execute_review
from backend_python.service.code_quality_service import Execute_enhance
from backend_python.schemas.dto.task import Task



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

