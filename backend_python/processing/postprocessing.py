# postprocess LLM output to json or
# file output


import uuid
import json
from fastapi.responses import FileResponse


async def postprocess(final_review: str):
    data = json.loads(final_review)
    return {
        "feedback": data["feedback"], 
        "issues": data.get("issues",[])
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
