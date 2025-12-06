from typing import Annotated, List, Dict, Any
from fastapi import APIRouter, UploadFile, File
from fastapi.responses import FileResponse
from backend_python.schema.context import CodeRequest, ReviewResponse
from backend_python.processing.postprocessing import postprocess, postprocess_md
from backend_python.processing.preprocessing import extract_chunks, process_uploaded_file
from backend_python.service.service import Execute
from pydantic import BaseModel

review_router = APIRouter(prefix="/analyse", tags=["analysis"])


@review_router.post("/code")
async def analyse_code(payload: CodeRequest) -> Dict[str, Any]:
    """
    Primary API endpoint for editor code analysis.

    Accepts editor code submit, processes the contents, and returns an analysis in json format.
    """

    chunked_context = extract_chunks(code=payload.submitted_code)
    print("")
    print(chunked_context)
    print("")
    response = await Execute(chunked_context)
    return postprocess(response)
   


@review_router.post("/export-md")
async def export_review(review: ReviewResponse) -> FileResponse:
    """
    Download analysis from the submitted editor code.

    """
    
    return postprocess_md(review)


# @review_router.post("/file")
# async def analyse_file(file: UploadFile) -> Dict[str, Any]:
#     """
#     Analyse python code from file, processes the contents, and returns an analysis in json format.

#     """

#     chunked_context = process_uploaded_file(file)
#     print("")
#     print(chunked_context)
#     print("")
#     response = await Execute(chunked_context) 
#     return postprocess(response)


# @review_router.post("/file/download")
# async def download_file_analysis(file: UploadFile) -> FileResponse:
#     """
#     Download analysis from an uploaded file of python code. 

#     """

#     chunked_context = process_uploaded_file(file)
#     print("")
#     print(chunked_context)
#     print("")
#     response = await Execute(chunked_context)
#     return postprocess_md(response)


@review_router.post("-multiple/")
async def upload_files(files: Annotated[List[UploadFile], File(description="Series of code files to be analysed by the LLM integration")] | None = None):
    # if len(files) == 0:
    #     return {"message": "failed to upload files"}
    # return {"filenames": [file.filename for file in files]}
    pass
