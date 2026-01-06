# import io 
# import csv
# import json
# from typing import Dict
# from backend_python.schemas.review.review_response import ReviewResponse


# def render_review_to_text(review: ReviewResponse) -> str:
#     lines = []
#     lines.append("# Code Review Summary\n")
#     lines.append(review.feedback + "\n")

#     if review.issues:
#         lines.append("## Issues\n")
#         for issue in review.issues:
#             lines.append(f"- **Line {issue.line}** ({issue.type}): {issue.description}")
    
#     return "\n".join(lines)


# def render_review_to_json(review: ReviewResponse) -> Dict[str,str]: 
#     json_obj = json.dumps(review.dict(), indent=2)
#     return json_obj

    
# def render_review_to_csv(review: ReviewResponse) -> str: 
#     output = io.StringIO()
#     writer = csv.writer(output)

#     writer.writerow(["Feedback", review.feedback])
#     writer.writerow([])

#     writer.writerow(["Line", "Type", "Description"])
#     for issue in review.issues:
#         writer.writerow([issue.line, issue.type, issue.description])

#     return output.getvalue()