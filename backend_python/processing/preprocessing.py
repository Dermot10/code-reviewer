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
import os
import ast
from typing import List
from tree_sitter import Parser
from tree_sitter_load import get_language_for_extension, NODE_TYPES
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


def extract_chunks(
    file_path: str, 
    code: str = "editor_input"
) -> List[CodeContext]:
    """
    Splits code into chunks by function or class.

    Args:
        code (str): Full code string from file/editor.
        file_path (str): Original file path (for metadata).

    Returns:
        List[CodeContext]: List of chunks ready for review.
    """

    # gets extension
    ext = file_path.split(".")[-1].lower()

    # maps extension to the language for the parser 
    language_obj, lang_name = get_language_for_extension(ext)

    # Fallback if unsupported language
    if not language_obj:
        return [CodeContext(file_path, "0", code, language=ext)]

    # set the language for the tree-sitter parser to use 
    
    parser = Parser()
    parser.set_language(language_obj)

    # parse file 
    try:
        tree = parser.parse(code.encode("utf8"))
    except Exception:
        return [CodeContext(file_path, "0", code, language=ext)]

    # convert source code into syntax tree 
    root = tree.root_node

    # the root node represents the whole file 
    # the child represents the nested constructs - classes, methods, funcs 
    target_nodes = NODE_TYPES.get(lang_name, [])
    chunks = []

    #recursively walk the tree, calling the func on child nodes 
    def walk(node):

        # if current node is class, method or func 
        if node.type in target_nodes:

            # are an attrubute of the treesitter node - a tuple
            # (line number, character position ), so [0] extracts the line number
            start = node.start_point[0]
            end = node.end_point[0]

            # extracts the extact lines of code that corresponding to that node
            block = "\n".join(code.splitlines()[start:end + 1])

            # Tree-sitter nodes don’t always have a text property, 
            # so the fallback ensures the chunk has an identifier.
            name = getattr(node, "text", f"node_{start}")

            # The chunk_id is unique per node in the file, based on its start line and type
            # e.g "3_function_definition"
            chunk_id = f"{start}_{node.type}"

            # wrap in code_context object
            # collect all matching nodes into a list
            chunks.append(
                CodeContext(
                    file_path=file_path,
                    chunk_id=chunk_id,
                    code=block,
                    language=ext,
                )
            )
        
        #recursively call on children - classes, methods, funcs
        for child in node.children:
            walk(child)

    # recursive call made 
    walk(root)

    # Fallback: no functions/classes → whole file returned
    if not chunks:
        chunks.append(CodeContext(file_path, "0", code, language=ext))


    # return list of chunks which correspond to the code snippets, 
    # to be given to the ai agent
    return chunks


