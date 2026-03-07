from backend_python.ai.ai_assistant_client import stream_open_ai_call
from backend_python.exceptions.exceptions import OpenAiProcessingError
from backend_python.schemas.ai.chat_context import ChatResponseContext
from backend_python.metrics import ASSISTANT_ERRORS
from backend_python.schemas.ai.prompts import ASSISTANT_SYSTEM_PROMPT

async def handle_assistant(
    prompt: str
) -> ChatResponseContext:
    """Streaming AI assistant agent"""

    system_prompt = ASSISTANT_SYSTEM_PROMPT

    try: 
        async for chunk in stream_open_ai_call(system_prompt, prompt): 
            yield chunk
    except OpenAiProcessingError as e: 
        ASSISTANT_ERRORS.inc()
        raise 