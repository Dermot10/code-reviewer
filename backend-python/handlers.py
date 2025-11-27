from typing import Annotated, List
from fastapi import APIRouter, UploadFile, File
from processing.preprocessing import chunk_code, process_uploaded_file
from service.service import Execute

review_router = APIRouter(prefix="/review", tags=["review"])


@review_router.post("/code")
async def analyse(submitted_code: str):
    """
    Primary API endpoint for file analysis.

    Accepts a file, processes the contents, and returns an analysis/summary.

    """

    download_flag = False
    chunked_code = chunk_code(submitted_code)
    response = await Execute(chunked_code, download_flag)
    # returns analysis from openai
    # return response


@review_router.post("/code/download")
async def analyse(submitted_code: str) -> File:
    """
    Primary API endpoint for file analysis.

    Accepts a file, processes the contents, and return a download file.

    """
    download_flag = True
    chunked_code = chunk_code(submitted_code)
    response = await Execute(chunked_code, download_flag)

    # returns analysis from openai

    # repackage into a file to return


@review_router.post("/file")
async def analyse(file: UploadFile):
    """
    Primary API endpoint for file analysis.

    Accepts a file, processes the contents, and returns an analysis/summary.

    """
    download_flag = False
    chunked_code = process_uploaded_file(file)
    # returns analysis from openai
    response = await Execute(chunked_code, download_flag)

    # return response


@review_router.post("/file/download")
async def analyse(file: UploadFile) -> File:
    """
    Primary API endpoint for file analysis.

    Accepts a file, processes the contents, and return a download file.

    """
    download_flag = True
    chunked_code = process_uploaded_file(file)
    # returns analysis from openai
    response = await Execute(chunked_code, download_flag)

    # repackage into a file


@review_router.post("-multiple/")
async def upload_files(files: Annotated[List[UploadFile], File(description="Series of code files to be analysed by the LLM integration")] | None = None):
    # if len(files) == 0:
    #     return {"message": "failed to upload files"}
    # return {"filenames": [file.filename for file in files]}
    pass
