import os
from typing import Type, Union
from openai import OpenAI
from dotenv import load_dotenv
from backend_python.schema.context import ReviewResponse, BestPracticesResponse
from backend_python.exceptions.exceptions import OpenAiProcessingError


load_dotenv()
api_key = os.getenv("OPENAI_API_KEY")


client = OpenAI()


def openai_call(system_prompt:str , input_prompt: str, output_format: Type[Union[ReviewResponse, BestPracticesResponse]]):
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
            text_format=output_format
        )

        review = response.output_parsed
    except Exception as e:
        raise OpenAiProcessingError(f"OpenAI API call failed: {str(e)}")

    return review