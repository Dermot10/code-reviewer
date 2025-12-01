# service for both file and non-file handling, flag can be passed down to signify whether to download file
# logging 
# metrics 
from typing import List
from fastapi import HTTPException
from backend_python.ai.ai_client import final_ai_call
from backend_python.processing.aggregator import aggregate_reviews
from backend_python.logger import logger
from backend_python.processing.context import CodeContext
from backend_python.ai.ai_agents import handle_syntax, handle_sematics, handle_best_practices, handle_security
from metrics import AGGREGATOR_ERRORS, AI_PROCESSING_TIME 


async def Execute(chunked_code: List[CodeContext]):
    try: 
        output = await code_review_service(chunked_code)
        return output
    except Exception as e: 
        logger.warning(f"failed to execute code review process - {e}")
        raise 

async def code_review_service(chunked_context: List[CodeContext]):
    chain = [
        handle_syntax,
        handle_sematics,
        handle_security,
        handle_best_practices
    ]

    results = []

    try: 
        with AI_PROCESSING_TIME.time():
            for agent in chain: 
                try: 
                    agent_result = await agent(chunked_context)
                    results.extend(agent_result)
                    logger.info(f"{agent.__name__} - succesfully processed code")
                except Exception as e: 
                    logger.error(f"{agent.__name__} failed: {e}")


            aggregated = aggregate_reviews(
                    code_contexts=chunked_context,
                    reviews=[results]
                )


            final_output = await final_ai_call(aggregated)

            return final_output 
    
    except Exception as e: 
        logger.warning(f"AI agents failed to process the code - {e}")
        raise HTTPException(status_code=500, detail="Internal server error")