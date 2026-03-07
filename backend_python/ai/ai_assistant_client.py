import os
from typing import Type, Union
from openai import AsyncOpenAI
from dotenv import load_dotenv
from backend_python.exceptions.exceptions import OpenAiProcessingError

load_dotenv()
api_key = os.getenv("OPENAI_API_KEY")


client = AsyncOpenAI()

async def stream_open_ai_call(system_prompt: str, prompt: str):
    stream = await client.chat.completions.create(
        model=os.getenv("MODEL"), 
        messages=[
            {"role": "system", "content": system_prompt}, 
            {"role": "user", "content": prompt},
        ], 
        stream=True
    )

    # incremental text from model, generated in streaming event
    # yielding data into memory for efficiency
    async for chunk in stream: 
        delta = chunk.choices[0].delta.content
        if delta: 
            yield delta