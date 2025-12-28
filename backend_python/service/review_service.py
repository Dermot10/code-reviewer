# service for both file and non-file handling, flag can be passed down to signify whether to download file
# logging 
# metrics 
from typing import List, Dict
from backend_python.logger import logger
from backend_python.schemas.ai.code_context import CodeContext
from backend_python.ai.ai_agents import handle_syntax
from backend_python.service.agent_service import agent_service ,aggregate_reviews


async def Execute_review(chunked_code: List[CodeContext]) -> Dict[str, str]:
    try: 
        output = await code_review_service(chunked_code)
        return output
    except Exception as e: 
        logger.warning(f"failed to execute code review service - {e}")
        raise 

async def code_review_service(chunked_context: List[CodeContext]) -> Dict[str, str]:
    chain = [
        handle_syntax,
    ]
    reviews = await agent_service(chain, chunked_context , aggregate_reviews)
    return reviews

