from pydantic import BaseModel

class ChatResponseContext(BaseModel):
    conversation_id: int
    chunk: str
    role: str = "assistant"