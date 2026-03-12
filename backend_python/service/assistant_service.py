from backend_python.logger import get_logger
from backend_python.ai.ai_assistant_agents import handle_assistant

logger = get_logger(__name__)


async def Execute_assistant(prompt: str): 
    try: 
        async for chunk in handle_assistant(prompt): 
            yield chunk
    except Exception as e: 
        logger.warning("failed to intiated assistant service", exc_info=e)
        raise 



