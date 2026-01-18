from typing import List, Dict
from fastapi import HTTPException
from backend_python.logger import get_logger
from backend_python.metrics import AGGREGATOR_ERRORS, AI_PROCESSING_TIME 
from backend_python.schemas.ai.code_context import CodeContext
from backend_python.schemas.ai.response_context import ResponseContext


logger = get_logger(__name__)

async def agent_service(chain: List, chunked_context: List[CodeContext], aggregate_func) -> Dict[str, str]:
    """
    Orchestrates multiple async agents to process chunked code and aggregates their outputs.

    This function handles:
        - Sending preprocessed code chunks to a chain of LLM-based agents.
        - Collecting individual agent outputs (ResponseContext objects).
        - Aggregating results using a provided aggregation strategy.
        - Handling per-agent errors without failing the entire pipeline.

    Args:
        chain (list): List of async agent functions to process chunked code.
        chunked_context (list[CodeContext]): Preprocessed code blocks to be analyzed.
        aggregate_func (Callable): Function to combine individual agent results into a final output.

    Returns:
        dict[str, str]: Aggregated output, in json like format.
    
    Raises:
        HTTPException: If aggregation or overall processing fails.
    """
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


def aggregate_reviews(results: List[ResponseContext]) -> Dict:
    """
    Aggregate code review results from multiple agents into a single summary.

    Processes the 'syntax', 'semantics', and 'security' fields of each ResponseContext,
    collecting structured issues and textual feedback.

    Args:
        results (list[ResponseContext]): List of agent review outputs.

    Returns:
        dict: Dictionary with:
            - 'feedback' (str): Concatenated textual feedback from all agents.
            - 'issues' (list[dict]): Structured list of issues with 'line', 'type', and 'description'.
    """
    issues = []
    feedback_lines = []

    for review in results:
        for field in ["syntax", "semantics", "security"]:
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

def aggregate_python_output(results: List[ResponseContext]) -> str:
    """
    Aggregate Python code output from multiple agents into a single module.

    Extracts the `.output` from the `best_practices` field of each ResponseContext and
    joins them with double newlines to produce a single, contiguous Python script.

    Args:
        results (list[ResponseContext]): List of agent outputs, each potentially containing a 'best_practices' field.

    Returns:
        str: Aggregated Python code as one module string.
    
    Notes:
        - Responses without 'best_practices' or with None are skipped safely.
        - The output preserves the order of results in the input list.
    """
    chunks = []
    for r in results:
        bp = getattr(r, "best_practices", None)
        if bp and getattr(bp, "output", None):
            chunks.append(bp.output)
    return "\n\n".join(chunks)
