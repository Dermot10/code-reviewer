# service for both file and non-file handling, flag can be passed down to signify whether to download file
# logging 
# metrics 
from typing import List
from fastapi import HTTPException
from backend_python.logger import logger
from backend_python.processing.context import CodeContext
from backend_python.ai.ai_agents import handle_syntax, handle_sematics, handle_best_practices, handle_security
from metrics import AGGREGATOR_ERRORS, AI_PROCESSING_TIME 


async def Execute(chunked_code: List[CodeContext], download_flag: bool):
    try: 
        if not download_flag:
            # replaced by func to wrap regular code review
            output = await handle_code_review(chunked_code)
            return output
        # func to wrap download functionality around code review
    except Exception as e: 
        logger.warning(f"failed to execute code review process - {e}")
        raise 

#func for downloadable file 


#func for regular 

# aggregator function for the final code review
async def handle_code_review(chunkedContext: List[CodeContext]):
    chain = [
        handle_syntax,
        handle_sematics,
        handle_best_practices,
        handle_security
    ]

    results = []

    try: 
        with AI_PROCESSING_TIME.time():
            for agent in chain: 
                try: 
                    result = await agent(chunkedContext)
                    results.append(result)
                    logger.info(f"{agent.__name__} - succesfully processed code")
                except Exception as e: 
                    logger.error(f"{agent.__name__} failed: {e}")
    
    except Exception as e: 
        logger.warning(f"AI agents failed to process the code - {e}")
        raise HTTPException(status_code=500, detail="Internal server error")