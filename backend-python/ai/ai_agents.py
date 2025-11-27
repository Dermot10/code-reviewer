from typing import List
from .ai_client import openai_call
from .prompts import SYNTAX_PROMPT, BEST_PRACTICES_PROMPT, SEMANTIC_PROMPT, SECURITY_PROMPT
from ..processing.preprocessing import CodeContext


async def handle_syntax(
    ctx: CodeContext
) -> None:
    prompt = f"{SYNTAX_PROMPT}\n\n{ctx.code}"
    response = await openai_call(prompt)
    return response


async def handle_best_practice(
    ctx: CodeContext
) -> None:
    prompt = f"{BEST_PRACTICES_PROMPT}\n\n{ctx.code}"
    response = await openai_call(prompt)
    return response


async def handle_sematics(
    ctx: CodeContext
) -> None:
    prompt = f"{SEMANTIC_PROMPT}\n\n{ctx.code}"
    response = await openai_call(prompt)
    return response


async def handle_security(
    ctx: CodeContext
) -> None:
    prompt = f"{SECURITY_PROMPT}\n\n{ctx.code}"
    response = await openai_call(prompt)
    return response


# aggregator function for the final code review
async def handle_code_review():
    ctx = CodeContext

    chain = [
        handle_code_review
    ]
