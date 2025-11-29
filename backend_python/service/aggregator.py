from typing import Dict, List
from backend_python.processing.preprocessing import CodeContext, ReviewContext

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
        if review.syntax is not None:
            final[review.chunk_id].syntax = review.syntax
        if review.semantics is not None:
            final[review.chunk_id].semantics = review.semantics
        if review.best_practices is not None:
            final[review.chunk_id].best_practices = review.best_practices
        if review.security is not None:
            final[review.chunk_id].security = review.security

    return final