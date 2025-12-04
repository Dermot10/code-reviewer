from typing import Annotated, List, Dict, Any
from fastapi import APIRouter, UploadFile, File
from fastapi.responses import FileResponse
from backend_python.processing.postprocessing import postprocess
from backend_python.processing.preprocessing import extract_chunks, process_uploaded_file
from backend_python.service.service import Execute
from pydantic import BaseModel

review_router = APIRouter(prefix="/analyse", tags=["analysis"])



class CodeRequest(BaseModel):
    submitted_code: str


@review_router.post("/code")
async def analyse(payload: CodeRequest) -> Dict[str, Any]:
    """
        Primary API endpoint for code analysis.

        Accepts editor code submit, processes the contents, and returns an analysis
        """

    chunked_context = extract_chunks(code=payload.submitted_code)
    print("")
    print(chunked_context)
    print("")
    response = await Execute(chunked_context)
    return postprocess(response)
   


@review_router.post("/code/download")
async def analyse(submitted_code: str) -> FileResponse:
    """
    Primary API endpoint for file analysis.

    Accepts editor code submit, processes the contents, and returns an analysis/summary as a downloadable file.

    """
 
    chunked_context = extract_chunks(submitted_code)
    response = await Execute(chunked_context)

    # returns analysis from openai

    # repackage into a file to return


@review_router.post("/file")
async def analyse(file: UploadFile) -> Dict[str, Any]:
    """
    Primary API endpoint for file analysis.

    Accepts a file, processes the contents, and returns an analysis/summary.

    """

    chunked_context = process_uploaded_file(file)
    # returns analysis from openai
    response = await Execute(chunked_context)

    # return response


@review_router.post("/file/download")
async def analyse(file: UploadFile) -> FileResponse:
    """
    Primary API endpoint for file analysis.

    Accepts a file, processes the contents, and returns an analysis/summary as a downloadable file.

    """

    chunked_context = process_uploaded_file(file)
    # returns analysis from openai
    response = await Execute(chunked_context)

    # repackage into a file


@review_router.post("-multiple/")
async def upload_files(files: Annotated[List[UploadFile], File(description="Series of code files to be analysed by the LLM integration")] | None = None):
    # if len(files) == 0:
    #     return {"message": "failed to upload files"}
    # return {"filenames": [file.filename for file in files]}
    pass
