from pydantic import BaseModel

class ResponseContext(BaseModel): 
    chunk_id: str
    syntax: str | None = None 
    semantics: str | None = None
    best_practices: str | None = None
    security: str | None = None
