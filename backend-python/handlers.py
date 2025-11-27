from typing import Annotated, List
from fastapi import APIRouter, UploadFile, File


router = APIRouter(prefix="/analyse", tags=["analyse"])


@router.post("/")
async def analyse(submitted_code: str):
    """
    Primary API endpoint for file analysis.

    Accepts a file, processes the contents, and returns an analysis/summary.

    """
    # response = await Execute(file)
    # returns analysis from openai
    # return response


@router.post("/")
async def analyse(submitted_code: str) -> File:
    """
    Primary API endpoint for file analysis.

    Accepts a file, processes the contents, and return a download file.

    """
    # response = await Execute(file)
    # returns analysis from openai
    # return response


@router.post("/")
async def analyse(file: UploadFile):
    """
    Primary API endpoint for file analysis.

    Accepts a file, processes the contents, and returns an analysis/summary.

    """
    # response = await Execute(file)
    # returns analysis from openai
    # return response


@router.post("/")
async def analyse(file: UploadFile) -> File:
    """
    Primary API endpoint for file analysis.

    Accepts a file, processes the contents, and return a download file.

    """
    # response = await Execute(file)
    # returns analysis from openai
    # return response


@router.post("-multiple/")
async def upload_files(files: Annotated[List[UploadFile], File(description="Series of code files to be analysed by the LLM integration")] | None = None):
    # if len(files) == 0:
    #     return {"message": "failed to upload files"}
    # return {"filenames": [file.filename for file in files]}
    pass
