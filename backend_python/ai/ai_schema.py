from langchain_core.prompts import PromptTemplate


review_prompt = PromptTemplate(
    input_variables=["code", "globals"],
    template="{system_prompt}\n{globals}\n\n# Code Chunk:\n{code}"
)

