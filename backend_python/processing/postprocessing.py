# postprocess LLM output to json or
# file output


def postprocess_review(final_review):
    return {
        "feedback": final_review["feedback"],
        "issues": final_review.get("issues", [])
    }

def postprocess_enhanced(enhanced_code):
    return {
        "enhanced_code": enhanced_code
    }


