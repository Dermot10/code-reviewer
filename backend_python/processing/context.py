from typing import List
from pydantic import BaseModel


class CodeContext(BaseModel):
    file_path: str
    chunk_id: str
    code: str
    # imports: List[str] = []
    # dependencies: List[str] = []
    # embedding_vector: List[float] = []


class ReviewContext(BaseModel): 
    chunk_id: str
    syntax: str | None = None 
    semantics: str | None = None
    best_practices: str | None = None
    security: str | None = None



