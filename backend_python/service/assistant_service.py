from typing import List, Dict
from fastapi import HTTPException
from backend_python.logger import get_logger
from backend_python.ai.ai_assistant_client import handle_assistant
from backend_python.schemas.ai.chat_context import ChatResponseContext


logger = get_logger(__name__)


async def Execute_assistant(prompt: str): 
    try: 
        async for chunk in handle_assistant(prompt): 
            yield chunk
    except Exception as e: 
        logger.warning("failed to intiated assistant service", exc_info=e)
        raise 



