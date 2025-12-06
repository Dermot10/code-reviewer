from typing import List, Dict, Any
from fastapi import HTTPException
from fastapi.responses import FileResponse
from backend_python.logger import logger
from backend_python.exceptions.exceptions import FileProcessingError
from backend_python.schema.context import ExportType, ReviewResponse
from backend_python.processing.fileprocessing import process_md, process_txt, process_csv, process_json
from backend_python.processing.renderer import render_review_to_text, render_review_to_json ,render_review_to_csv


async def Exceute_export(export_type: str, review: ReviewResponse) -> FileResponse:
    try: 
        response_file = await export_file_service(export_type, review)
        return response_file
    except Exception as e: 
        logger.warning(f"failed to execute export file service - {e}")
        raise 


def export_file_service(export_type: str, review: ReviewResponse) -> FileResponse: 
    try: 
        if export_type == ExportType.MD:
            final_review = render_review_to_text(review)
            file_response = process_md(export_type, final_review)
            print("Exporting as .md file")
            return file_response

        elif export_type == ExportType.TXT:
            final_review = render_review_to_text(review)
            file_response = process_txt(export_type, review)
            print("Exporting as .txt file")
            return file_response

        elif export_type == ExportType.CSV:
            final_review = render_review_to_csv(review)
            file_response = process_csv(export_type, review)
            print("Exporting as .csv file")
            return file_response

        elif export_type == ExportType.JSON:
            final_review = render_review_to_json(review)
            file_response = process_json(export_type, review)
            print("Exporting as .json file")
            return file_response
        
        # elif export_type == ExportType.PY:
        #     file_response = process_py(export_type, review)
        #     print("Exporting as .py file")
        #     return file_response
        else: 
            print(f"unknown export typ {export_type}")
            return

    except FileProcessingError as e: 
        logger.warning(f"failed to export file - {e}")
        raise HTTPException(status_code=500, detail="Internal server error")