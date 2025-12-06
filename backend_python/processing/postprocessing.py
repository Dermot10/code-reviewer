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


