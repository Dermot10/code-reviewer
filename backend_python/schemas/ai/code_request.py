from pydantic import BaseModel

class CodeRequest(BaseModel):
    submitted_code: str