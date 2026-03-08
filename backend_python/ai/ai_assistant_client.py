import os
from openai import AsyncOpenAI
from dotenv import load_dotenv, find_dotenv
from backend_python.exceptions.exceptions import OpenAiProcessingError

load_dotenv(find_dotenv())

client = AsyncOpenAI(api_key=os.getenv("OPENAI_API_KEY"))

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