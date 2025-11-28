import os
import asyncio
from openai import OpenAI
from dotenv import load_dotenv

from exceptions.exceptions import OpenAiProcessingError
from prompts import SYSTEM_PROMPT

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
        response = client.responses.create(
            model=os.getenv("MODEL"),
            instructions=SYSTEM_PROMPT,
            input=input_prompt
        )
    except Exception as e:
        raise OpenAiProcessingError(f"OpenAI API call failed: {str(e)}")

    return response.output_text
