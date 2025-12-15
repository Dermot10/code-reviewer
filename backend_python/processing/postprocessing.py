# postprocess LLM output to json or
# file output


def postprocess_review(final_review):
    return {
        "feedback": final_review["feedback"],
        "issues": final_review.get("issues", [])
    }

def postprocess_enhanced(enhanced_code):
    enhanced_output = {
        "enhanced_code": enhanced_code
    }
    print(f"this is the enhanced output sent to the front end - {enhanced_output}")
    return enhanced_output

