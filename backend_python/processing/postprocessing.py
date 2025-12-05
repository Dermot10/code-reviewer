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

def postprocess_md(final_review):
    file_path = f"/tmp/review_{uuid.uuid4()}.md"

    with open(file_path, "w", encoding="utf-8") as f:
        f.write(final_review)

    return FileResponse(
        path=file_path,
        media_type="text/markdown",
        filename="code_review.md"
    )
