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
from backend_python.processing.context import CodeContext
from backend_python.metrics import CHUNKING_FAILURES, FILE_EXTRACTION_ERRORS


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
        FILE_EXTRACTION_ERRORS.inc()
        raise ValueError(f"Failed to read file {file_path}: {str(e)}")

    # Use your existing chunking function
    chunks = extract_chunks(code_body, file_path=file_path)
    return chunks

def extract_chunks(file_path: str, code: str = "editor_input") -> List[CodeContext]:
    ext = file_path.split(".")[-1].lower()

    # Only parse Python code; fallback for other languages
    if ext != "py":
        return [CodeContext(file_path, "0", code, language=ext)]

    chunks = []

    try:
        tree = ast.parse(code)
    except Exception:
        return [CodeContext(file_path, "0", code, language=ext)]

    for node in ast.walk(tree):
        if isinstance(node, (ast.FunctionDef, ast.AsyncFunctionDef, ast.ClassDef)):
            #file line starts at 1, strings start at 0
            start = node.lineno - 1
            end = node.end_lineno - 1 if hasattr(node, "end_lineno") else start
            block = "\n".join(code.splitlines()[start:end + 1])

            # Use node type and name for chunk ID if available
            name = getattr(node, "name", f"node_{start}")
            chunk_id = f"{start}_{type(node).__name__}_{name}"

            chunks.append(
                CodeContext(
                    file_path=file_path,
                    chunk_id=chunk_id,
                    code=block,
                    ext=ext
                )
            )

    # Fallback: whole file if no classes/functions
    if not chunks:
        chunks.append(CodeContext(file_path, "0", code, language=ext))

    return chunks