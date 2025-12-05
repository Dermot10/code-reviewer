from typing import Dict, List
from backend_python.schema.context import CodeContext, ReviewContext

def aggregate_reviews(
    code_contexts: List[CodeContext], 
    reviews: List[List[ReviewContext]]
) -> Dict[str, ReviewContext]: 

    final = {}

    #flatten
    merged = [item for sublist in reviews for item in sublist]

    for context in code_contexts: 
        final[context.chunk_id] = ReviewContext(chunk_id=context.chunk_id)

    
    # merge agents → combined ReviewContext
    for review in merged:
        target = final[review.chunk_id]
        for field in ReviewContext.model_fields:
            value = getattr(review, field)
            if value is not None:
                # add the agent's specific work to the final 
                setattr(target, field, value)

    return final
