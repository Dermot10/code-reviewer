from pydantic import BaseModel
from typing import Union, Literal


class ReviewTask(BaseModel):
    type: Literal["review"] 
    user_id: int 
    review_id:int
    action: str = "generate_summary"


class EnhancementTask(BaseModel):
    type: Literal["review"]
    user_id: int
    enhancement_id: int 
    action: str = "enhance_code"


Task = Union[ReviewTask, EnhancementTask]