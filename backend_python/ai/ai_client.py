import os
import asyncio
from openai import OpenAI
from dotenv import load_dotenv
from typing import Dict
from backend_python.schema.context import ReviewResponse
from backend_python.schema.context import ReviewContext
from backend_python.exceptions.exceptions import OpenAiProcessingError


load_dotenv()
api_key = os.getenv("OPENAI_API_KEY")


client = OpenAI()


def openai_call(system_prompt ,input_prompt: str):
    """
    OpenAI API call

    Accepts extracted file data and sends it to LLM for generative response. 

    """
    try:
        # must be synchronous, to allow LLM to complete text generation
        response = client.responses.parse(
            model=os.getenv("MODEL"),
            input=[
            {"role": "system", "content": system_prompt},
            {"role": "user", "content": input_prompt},
            ],
            text_format=ReviewResponse
        )

        review = response.output_parsed
    except Exception as e:
        raise OpenAiProcessingError(f"OpenAI API call failed: {str(e)}")

    return review