import uuid
import json
from typing import Dict
from fastapi.responses import FileResponse





def create_file_path(export_type: str, content: str) -> str:
    file_path = f"/tmp/exported_review.{export_type.lower()}"
    with open(file_path, "w", encoding="utf-8") as f:
        f.write(content)
    return file_path

def process_md(export_type, final_review):
    file_path = create_file_path(export_type, final_review)

    return FileResponse(
        path=file_path,
        media_type="text/markdown",
        filename="code_review.md"
    )

def process_txt(export_type: str, final_review):
    file_path = create_file_path(export_type, final_review)

    return FileResponse(
        path=file_path,
        media_type="text/plain",
        filename="code_review.txt"
    )


def process_json(export_type, final_review):
    file_path = create_file_path(export_type, final_review)

    return FileResponse(
        path=file_path,
        media_type="application/json",
        filename="code_review.json"
    )

def process_csv(export_type, final_review):
    file_path = create_file_path(export_type, final_review)

    return FileResponse(
        path=file_path,
        media_type="text",
        filename="code_review.csv"
    )