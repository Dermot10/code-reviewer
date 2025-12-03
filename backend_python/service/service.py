# service for both file and non-file handling, flag can be passed down to signify whether to download file
# logging 
# metrics 
from typing import List, Dict, Any
from fastapi import HTTPException
from backend_python.ai.ai_client import final_ai_call
from backend_python.processing.aggregator import aggregate_reviews
from backend_python.logger import logger
from backend_python.processing.context import CodeContext, ReviewContext
from backend_python.ai.ai_agents import handle_syntax, handle_sematics, handle_best_practices, handle_security
from backend_python.metrics import AGGREGATOR_ERRORS, AI_PROCESSING_TIME 


async def Execute(chunked_code: List[CodeContext]) -> Dict[str, str]:
    try: 
        output = await code_review_service(chunked_code)
        return output
    except Exception as e: 
        logger.warning(f"failed to execute code review process - {e}")
        raise 

async def code_review_service(chunked_context: List[CodeContext]) -> Dict[str, str]:
    chain = [
        handle_syntax,
        # handle_sematics,
        # handle_security,
        # handle_best_practices
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
            print("")
            print(f"The extended results from the agent call ---->> {results}")
            print("")
        aggregated_results = aggregate_to_preprocess(results)
        print(f"The aggregated results to be postprocessed ---->> {aggregated_results}")
        
        return aggregated_results
    
    except Exception as e: 
        logger.warning(f"AI agents failed to process the code - {e}")
        raise HTTPException(status_code=500, detail="Internal server error")

def aggregate_to_preprocess(results: List[ReviewContext]) -> Dict:
    issues = []
    feedback_lines = []

    for review in results:
        for field in ["syntax", "semantics", "best_practices", "security"]:
            val = getattr(review, field)
            if val:
                # collect issues
                if hasattr(val, "issues") and val.issues:
                    issues.extend(
                        {"line": i.line, "type": i.type, "description": i.description}
                        for i in val.issues
                    )
                # collect feedback text
                if hasattr(val, "feedback") and val.feedback:
                    feedback_lines.append(val.feedback)

    # combine all feedbacks into a single concise string
    combined_feedback = " | ".join(feedback_lines) if feedback_lines else "No issues found."

    return {"feedback": combined_feedback, "issues": issues}


# final report logic
# aggregated = aggregate_reviews(
            #         code_contexts=chunked_context,
            #         reviews=[results]
            #     )