# Purpose: Runs linters, type checkers, style formatters, and maybe even small “improve code” comments.

# Returns: rewritten/improved code snippet.

from typing import List, Dict
from backend_python.logger import logger
from backend_python.schemas.ai.code_context import CodeContext
from backend_python.ai.ai_agents import handle_best_practices
from backend_python.service.agent_service import agent_service, aggregate_python_output


async def Execute_enhance(chunked_code: List[CodeContext]) -> Dict[str, str]: 
    try: 
        output = await code_quality_service(chunked_code)
        return output
    except Exception as e:
        logger.warning(f"failed to initiate enhance code process - {e}")
        raise 



async def code_quality_service(chunked_context: List[CodeContext]) -> Dict[str, str]: 

    chain = [
        handle_best_practices
    ]

    enhanced_code = await agent_service(chain, chunked_context, aggregate_python_output)

    return enhanced_code

