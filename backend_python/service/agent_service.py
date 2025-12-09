from typing import List, Dict
from fastapi import HTTPException
from backend_python.logger import logger
from backend_python.metrics import AGGREGATOR_ERRORS, AI_PROCESSING_TIME 
from backend_python.schema.context import CodeContext, ReviewContext


async def agent_service(chain: List, chunked_context: List[CodeContext], aggregate_func) -> Dict[str, str]:

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
            print(f"The extended results from the agent call ---->> \n{results}")
            print("")
        aggregated_results = aggregate_func(results)
        print(f"The aggregated results to be postprocessed ---->> \n{aggregated_results}")
        
        return aggregated_results
    
    except Exception as e: 
        logger.warning(f"AI agents failed to process the code - {e}")
        raise HTTPException(status_code=500, detail="Internal server error")


def aggregate_reviews(results: List[ReviewContext]) -> Dict:
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

    
    combined_feedback = "\n\n".join(feedback_lines)

    return {"feedback": combined_feedback, "issues": issues}

def aggregate_python_output(results: List[ReviewContext]) -> str:
    chunks = [getattr(r, "best_practices").output
                for r in results
                if getattr(r, "best_practices", None)]
    return "\n\n".join(chunks)
