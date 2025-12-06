from typing import List, Dict, Any
from fastapi import HTTPException
from backend_python.logger import logger
from backend_python.exceptions.exceptions import FileProcessingError
from backend_python.schema.context import ExportType, ReviewResponse
from backend_python.processing.postprocessing import process_md, process_txt


async def Exceute_export(export_type: str, review: ReviewResponse):
    try: 
        output = await export_file_service(export_type)
        return output
    except Exception as e: 
        logger.warning(f"failed to execute export file service - {e}")
        raise 


def export_file_service(export_type: str, review: ReviewResponse): 
    try: 
        if export_type == ExportType.MD:
            file_response = process_md(export_type, review)
            print("Exporting as .md file")
            return file_response
        elif export_type == ExportType.TXT:
            file_response = process_txt(export_type, review)
            print("Exporting as .txt file")
            return file_response
        elif export_type == ExportType.CSV:
            file_response = process_csv(export_type, review)
            print("Exporting as .csv file")
            return file_response
        elif export_type == ExportType.JSON:
            file_response = process_json(export_type, review)
            print("Exporting as .json file")
            return file_response
        elif export_type == ExportType.PY:
            file_response = process_py(export_type, review)
            print("Exporting as .py file")
            return file_response
        else: 
            print(f"unknown export typ {export_type}")
            return

    except FileProcessingError as e: 
        logger.warning(f"failed to export file - {e}")
        raise HTTPException(status_code=500, detail="Internal server error")