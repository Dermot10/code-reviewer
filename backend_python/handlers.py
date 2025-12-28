from typing import Annotated, List, Dict, Any
from fastapi import APIRouter, UploadFile, File
from fastapi.responses import FileResponse
from backend_python.schemas.ai.code_request import CodeRequest
from backend_python.schemas.review.review_response import ReviewResponse
from backend_python.schemas.common.enums import ExportType
from backend_python.processing.postprocessing import postprocess_review, postprocess_enhanced
from backend_python.processing.preprocessing import extract_chunks
from backend_python.service.export_service import Exceute_export
from backend_python.service.review_service import Execute_review
from backend_python.service.code_quality_service import Execute_enhance


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
    response = await Execute_review(chunked_context)
    return postprocess_review(response)
   


@review_router.post("/code-quality")
async def enhance_code(payload: CodeRequest) -> Dict[str, Any]:
    """
    Enhances the quality of the submitted code with AI suggestion.

    Accepts editor code submit, processes the contents, and returns the code improved or better optimised.

    """
    chunked_context = extract_chunks(code=payload.submitted_code)
    print("")
    print(chunked_context)
    print("")
    response = await Execute_enhance(chunked_context)
    print("----- The response from the code quality -----")
    print(response)
    print("")
    return postprocess_enhanced(response)

@review_router.post("/export-md")
async def export_review_md(review: ReviewResponse) -> FileResponse:
    """
    Download analysis from the submitted editor code.

    """
    export_choice = ExportType.MD.value
    print(export_choice)
    export_file = Exceute_export(export_choice, review)
    return export_file


@review_router.post("/export-json")
async def export_review_json(review: ReviewResponse) -> FileResponse:
    """
    Download analysis as json file.

    """
    export_choice = ExportType.JSON.value
    export_file = Exceute_export(export_choice, review)
    return export_file
    

@review_router.post("/export-csv")
async def export_review_csv(review: ReviewResponse) -> FileResponse:
    """
    Download analysis as csv file.

    """
    export_choice = ExportType.CSV.value
    export_file = Exceute_export(export_choice, review)
    return export_file
    

@review_router.post("/export-txt")
async def export_review_txt(review: ReviewResponse) -> FileResponse:
    """
    Download analysis as txt file.

    """
    export_choice = ExportType.TXT.value
    export_file = Exceute_export(export_choice, review)
    return export_file

# @review_router.post("-multiple/")
# async def upload_files(files: Annotated[List[UploadFile], File(description="Series of code files to be analysed by the LLM integration")] | None = None):
#     # if len(files) == 0:
#     #     return {"message": "failed to upload files"}
#     # return {"filenames": [file.filename for file in files]}
#     pass
