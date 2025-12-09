from typing import List
from pydantic import BaseModel
from enum import Enum


class ExportType(Enum): 
    MD = "markdown"
    TXT = "txt"
    CSV = "csv"
    JSON = "json"
    PY = "py"

class CodeRequest(BaseModel):
    submitted_code: str

class CodeContext(BaseModel):
    file_path: str
    chunk_id: str
    code: str
    globals: List[str] = []
    # imports: List[str] = []
    # dependencies: List[str] = []
    # embedding_vector: List[float] = []

class ReviewContext(BaseModel): 
    chunk_id: str
    syntax: str | None = None 
    semantics: str | None = None
    best_practices: str | None = None
    security: str | None = None

class Issue(BaseModel): 
    line: int 
    type: str # bug|security|style|other
    description: str 

class ReviewResponse(BaseModel): 
    feedback: str
    issues: List[Issue]


class BestPracticesResponse(BaseModel): 
    output: str


