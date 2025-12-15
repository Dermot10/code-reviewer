# preprocess the incoming request
# chunking large files
# chunking funcs
# preserve imports, relationship
# relevant points to LLM


import ast
from uuid import uuid4
from typing import List, Optional
from backend_python.schemas.ai.code_context import CodeContext
from backend_python.metrics import CHUNKING_FAILURES, FILE_EXTRACTION_ERRORS


def extract_globals(tree: ast.AST, code_lines: List[str]) -> List[str]:
    globals_found = []

    for node in tree.body:
        
        # Handle simple assignment: x = 10
      
        if isinstance(node, ast.Assign):
            for target in node.targets:
                if isinstance(target, ast.Name):
                    # Pull the actual source exactly as written in the file
                    line = code_lines[node.lineno - 1].rstrip()
                    globals_found.append(line)

       
        # Annotated assignment: x: int = 10
      
        elif isinstance(node, ast.AnnAssign):
            if isinstance(node.target, ast.Name):
                line = code_lines[node.lineno - 1].rstrip()
                globals_found.append(line)

   
        # The __all__ variable
       
        if (
            isinstance(node, (ast.Assign, ast.AnnAssign))
            and any(
                isinstance(t, ast.Name) and t.id == "__all__"
                for t in (node.targets if isinstance(node, ast.Assign) else [node.target])
            )
        ):
            line = code_lines[node.lineno - 1].rstrip()
            globals_found.append(line)

       
        # Imports count as global declarations
       
        if isinstance(node, (ast.Import, ast.ImportFrom)):
            line = code_lines[node.lineno - 1].rstrip()
            globals_found.append(line)

    return globals_found


def extract_chunks(code: str, file_path: Optional[str]= None) -> List[CodeContext]:
    """
    Extract code chunks for processing.

    Args:
        code: The Python code to process.
        file_path: Optional path of the file the code came from.

    Returns:
        List[CodeContext]: One chunk per function/class, a chunk being a code block
    """
    chunks = []
    path = file_path or ""

    try:
        tree = ast.parse(code) 
    except Exception:
        return [
            CodeContext (
            file_path=path, 
            chunk_id=str(uuid4()),
            code=code,
            globals= [],
            )
        ]
    lines = code.splitlines() 
    globals_list = extract_globals(tree, lines)

    for node in ast.walk(tree):
        if isinstance(node, (ast.FunctionDef, ast.AsyncFunctionDef, ast.ClassDef)):
            # file line starts at 1, strings start at 0
            start = node.lineno - 1
            end = node.end_lineno - 1 if hasattr(node, "end_lineno") else start
            block = "\n".join(lines[start:end + 1])

            # Use node type and name for chunk ID if available
            name = getattr(node, "name", f"node_{start}")
            chunk_id = f"{start}_{type(node).__name__}_{name}"

            chunks.append(
                CodeContext(
                    file_path=path,
                    chunk_id=chunk_id,
                    code=block,
                    globals=globals_list,
                )
            )

    # Fallback: whole file if no classes/functions
    if not chunks:
        chunks.append(CodeContext(file_path=path, chunk_id="0", code=code, globals=[]))

    return chunks


def process_uploaded_file(file_path: str) -> List[CodeContext]:
    """
    Reads a python file and returns a list of CodeContext chunks.

    Args:
        file_path (str): Path to the uploaded file.

    Returns:
        List[CodeContext]: One chunk per function/class
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