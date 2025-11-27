# preprocess the incoming request
# chunking large files
# chunking funcs
# preserve imports, relationship
# relevant points to LLM

# segment for LLM agents, chain of responsibility pattern
# Syntax & structure pass → highlight obvious issues
# Semantic analysis pass → logic flaws, bugs, confusion
# Best practices pass → style, clarity, naming
# Security pass → SQL injection, unsafe file handling, etc


# Language support
# add lang specific rules
# Use language server protocols (LSP) later if you want deeper static analysis

import ast
from typing import List
from .context import CodeContext


def process_uploaded_file(file_path: str) -> List[CodeContext]:
    """
    Reads a code file and returns a list of CodeContext chunks.

    Args:
        file_path (str): Path to the uploaded file.

    Returns:
        List[CodeContext]: One chunk per function/class (or whole file fallback)
    """
    try:
        with open(file_path, "r", encoding="utf-8") as f:
            code_body = f.read()
    except Exception as e:
        raise ValueError(f"Failed to read file {file_path}: {str(e)}")

    # Use your existing chunking function
    chunks = chunk_code(code_body, file_path=file_path)
    return chunks


def chunk_code(
    code: str,
    file_path: str = "editor_input"
) -> List[CodeContext]:
    """
    Splits code into chunks by function or class.

    Args:
        code (str): Full code string from file/editor.
        file_path (str): Original file path (for metadata).

    Returns:
        List[CodeContext]: List of chunks ready for review.
    """
    chunks: List[CodeContext] = []
    try:
        tree = ast.parse(code)
    except SyntaxError:
        # fallback: treat entire file as one chunk if parse fails
        return [CodeContext(file_path=file_path, chunk_id="0", code=code)]

    for i, node in enumerate(tree.body):
        if isinstance(node, (ast.FunctionDef, ast.ClassDef)):
            start = node.lineno - 1
            end = getattr(node, "end_lineno", start + 1)
            chunk_code = "\n".join(code.splitlines()[start:end])
            chunk_id = f"{i}_{node.name}"
            chunks.append(
                CodeContext(
                    file_path=file_path,
                    chunk_id=chunk_id,
                    code=chunk_code
                )
            )

    # fallback: if no functions/classes found, use entire file
    if not chunks:
        chunks.append(CodeContext(
            file_path=file_path, chunk_id="0", code=code))

    return chunks
