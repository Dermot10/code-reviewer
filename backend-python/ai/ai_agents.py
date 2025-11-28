from typing import List

from .ai_client import openai_call
from .prompts import SYNTAX_PROMPT, BEST_PRACTICES_PROMPT, SEMANTIC_PROMPT, SECURITY_PROMPT
from ..processing.preprocessing import CodeContext, ReviewContext


async def handle_agent(
    codeContexts: List[CodeContext], 
    prompt: str, 
    strategy: str
) -> List[ReviewContext]:
    results = []
    for context in codeContexts:
        ai_prompt = f"{prompt}\n\n{context.code}"
        review = ReviewContext(chunk_id = context.chunk_id)
        setattr(review, strategy,await openai_call(ai_prompt)),
        results.append(review)
    return results

async def handle_syntax(
    codeContexts: List[CodeContext]
) -> List[ReviewContext]:
   return await handle_agent(codeContexts, SYNTAX_PROMPT, "syntax")

async def handle_best_practice(
    codeContexts: List[CodeContext]
) -> List[ReviewContext]:
   return await handle_agent(codeContexts, BEST_PRACTICES_PROMPT, "best_practices")

async def handle_sematics(
    codeContexts: List[CodeContext]
) -> List[ReviewContext]:
   return await handle_agent(codeContexts, SEMANTIC_PROMPT, "semantics")

async def handle_security(
    codeContexts: List[CodeContext]
) -> List[ReviewContext]:
   return await handle_agent(codeContexts, SECURITY_PROMPT, "syntax")


# aggregator function for the final code review
async def handle_code_review():
    ctx = CodeContext

    chain = [
        handle_code_review
    ]
