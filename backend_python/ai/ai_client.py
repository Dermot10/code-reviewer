import os
import asyncio
from openai import OpenAI
from dotenv import load_dotenv
from typing import Dict
from backend_python.schema.context import ReviewResponse
from backend_python.schema.context import ReviewContext
from backend_python.exceptions.exceptions import OpenAiProcessingError
from backend_python.ai.prompts import SYSTEM_PROMPT

load_dotenv()
api_key = os.getenv("OPENAI_API_KEY")


client = OpenAI()


def openai_call(input_prompt: str):
    """
    OpenAI API call

    Accepts extracted file data and sends it to LLM for generative response. 

    """
    try:
        # must be synchronous, to allow LLM to complete text generation
        response = client.responses.parse(
            model=os.getenv("MODEL"),
            input=[
            {"role": "system", "content": SYSTEM_PROMPT},
            {"role": "user", "content": input_prompt},
            ],
            text_format=ReviewResponse
        )

        review = response.output_parsed
    except Exception as e:
        raise OpenAiProcessingError(f"OpenAI API call failed: {str(e)}")

    return review



async def final_ai_call(aggregated: Dict[str, ReviewContext]):
    """
    Combines all reviews and calls the LLM for a final coherent summary.
    """

    # Build a total report
    report = ""
    for chunk_id, review in aggregated.items():
        report += f"\n### Chunk {chunk_id}\n"

        if review.syntax:
            report += f"\n**Syntax Review:**\n{review.syntax}\n"
        if review.semantics:
            report += f"\n**Semantics Review:**\n{review.semantics}\n"
        if review.best_practices:
            report += f"\n**Best Practices:**\n{review.best_practices}\n"
        if review.security:
            report += f"\n**Security:**\n{review.security}\n"

    final_prompt = f"""
You are an expert python code reviewer.

You will now create a final, unified code review based on the following chunk-level agent feedback:

{report}

Return a **single coherent, structured review** with:
- A high-level summary
- Per-chunk conclusions
- Overall assessment
- Suggested rewritten code sections (if needed)
"""

    polished = openai_call(final_prompt)

    return {
        "chunks": aggregated,
        "final_review": polished
    }
