# postprocess LLM output to json or
# file output


import uuid
import json
from typing import Dict
from fastapi.responses import FileResponse


def postprocess_review(final_review):
    return {
        "feedback": final_review["feedback"],
        "issues": final_review.get("issues", [])
    }

def postprocess_enhanced(enhanced_code):
    return {
        "enhanced_code": ""
    }



def create_file_path(export_type ,final_review): 
    ext = export_type.lower()
    file_path = f"/tmp/review_{uuid.uuid4()}.{ext}"

    with open(file_path, "w", encoding="utf-8") as f:
        f.write(final_review)

    return file_path

def process_md(export_type, final_review):
    file_path = create_file_path(export_type, final_review)

    return FileResponse(
        path=file_path,
        media_type="text/markdown",
        filename="code_review.md"
    )

def process_txt(export_type, final_review):
    file_path = create_file_path(export_type, final_review)

    return FileResponse(
        path=file_path,
        media_type="text",
        filename="code_review.txt"
    )


# other file processors 
