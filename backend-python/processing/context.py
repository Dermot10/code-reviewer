from typing import List
from pydantic import BaseModel


class CodeContext(BaseModel):
    file_path: str
    chunk_id: str
    code: str
    imports: List[str] = []
    dependencies: List[str] = []
    embedding_vector: List[float] = []
