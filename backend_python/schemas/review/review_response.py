from typing import List
from pydantic import BaseModel
from backend_python.schemas.review.issue import Issue

class ReviewResponse(BaseModel): 
    feedback: str
    issues: List[Issue]
