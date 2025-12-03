from typing import List

from backend_python.ai.ai_client import openai_call
from backend_python.ai.prompts import SYNTAX_PROMPT, BEST_PRACTICES_PROMPT, SEMANTIC_PROMPT, SECURITY_PROMPT
from backend_python.processing.context import CodeContext, ReviewContext
from backend_python.metrics import SYNTAX_ERRORS, SEMANTICS_ERRORS, BEST_PRACTICES_ERRORS, SECURITY_ERRORS
from backend_python.exceptions.exceptions import OpenAiProcessingError


async def handle_syntax(
    code_contexts: List[CodeContext]
) -> List[ReviewContext]:
    try:
        return await handle_agent(code_contexts, SYNTAX_PROMPT, "syntax")
    except OpenAiProcessingError as e: 
        SYNTAX_ERRORS.inc()
        raise
        
async def handle_sematics(
    code_contexts: List[CodeContext]
) -> List[ReviewContext]:
    try:
        return await handle_agent(code_contexts, SEMANTIC_PROMPT, "semantics")
    except OpenAiProcessingError as e: 
        SEMANTICS_ERRORS.inc()
        raise

async def handle_best_practices(
    code_contexts: List[CodeContext]
) -> List[ReviewContext]:
    try: 
        return await handle_agent(code_contexts, BEST_PRACTICES_PROMPT, "best_practices")
    except OpenAiProcessingError as e: 
        BEST_PRACTICES_ERRORS.inc()
        raise

async def handle_security(
    code_contexts: List[CodeContext]
) -> List[ReviewContext]:
    try:
        return await handle_agent(code_contexts, SECURITY_PROMPT, "syntax")
    except OpenAiProcessingError as e: 
        SECURITY_ERRORS.inc()
        raise

async def handle_agent(
    code_contexts: List[CodeContext],
    prompt: str,
    strategy: str
) -> List[ReviewContext]:
    results = []
    for context in code_contexts:
        ai_prompt = f"{prompt}\n\n{context.code}"
        review = ReviewContext(chunk_id = context.chunk_id)
        setattr(review, strategy, await openai_call(ai_prompt)),
        results.append(review) # may extend into one final result
    return results