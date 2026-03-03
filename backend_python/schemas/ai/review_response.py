from typing import List
from pydantic import BaseModel


class BestPracticesResponse(BaseModel): 
    output: str


class Issue(BaseModel): 
    line: int 
    type: str # bug|security|style|other
    description: str 


class ReviewResponse(BaseModel): 
    feedback: str
    issues: List[Issue]
