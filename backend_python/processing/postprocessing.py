# postprocess LLM output to json or
# file output


import uuid
from fastapi.responses import FileResponse


async def postprocess(final_review: str, download: bool):
    if not download:
        return {"review": final_review}


async def postprocess(final_review: str, download: bool):
    if download:
        file_path = f"/tmp/review_{uuid.uuid4()}.txt"
        with open(file_path, "w") as f:
            f.write(final_review)

        return FileResponse(
            file_path,
            media_type="text/plain",
            filename="code_review.txt"
        )
