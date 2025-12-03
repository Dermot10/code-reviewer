from pydantic import BaseModel
from typing import List 


class Issue(BaseModel): 
    line: int 
    type: str # bug|security|style|other
    description: str 


class ReviewResponse(BaseModel): 
    feedback: str
    issues: List[Issue]

