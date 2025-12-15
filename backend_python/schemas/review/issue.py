from pydantic import BaseModel


class Issue(BaseModel): 
    line: int 
    type: str # bug|security|style|other
    description: str 