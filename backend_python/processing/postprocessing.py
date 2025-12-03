# postprocess LLM output to json or
# file output


import uuid
import json
from typing import Dict
from fastapi.responses import FileResponse


def postprocess(final_review):
    return {
        "feedback": final_review["feedback"],
        "issues": final_review.get("issues", [])
    }

async def postprocess_file(final_review: str):
    file_path = f"/tmp/review_{uuid.uuid4()}.txt"
    with open(file_path, "w") as f:
        f.write(final_review)

    return FileResponse(
        file_path,
        media_type="text/plain",
        filename="code_review.txt"
    )


#extension for select langauages to create the correct file type 
