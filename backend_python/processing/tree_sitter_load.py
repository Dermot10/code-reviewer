from pathlib import Path
from tree_sitter import Language, Parser
from typing import Tuple
import os

LANG_LIB_PATH = os.path.join(os.path.dirname(__file__), "languages.so")

# Map file extensions -> Tree-sitter Language ID
EXTENSION_TO_LANGUAGE = {
    "py": "python",
    "js": "javascript",
    "ts": "typescript",
    "go": "go",
    "java": "java",
    "cs": "c_sharp",
    "rb": "ruby",
    "php": "php",
    "rs": "rust",
}

# Supported node types per language
NODE_TYPES = {
    "python": ["function_definition", "class_definition"],
    "javascript": ["function_declaration", "method_definition", "class_declaration"],
    "typescript": ["function_declaration", "method_definition", "class_declaration"],
    "go": ["function_declaration", "method_declaration"],
    "java": ["class_declaration", "method_declaration"],
    "c_sharp": ["class_declaration", "method_declaration"],
    "php": ["function_definition", "method_declaration", "class_declaration"],
    "ruby": ["method", "class"],
    "rust": ["function_item", "impl_item"],
}

# Load the combined library
LANGUAGES = {
    # dict comp for language name: Language func(language_path, language_name)
    lang_name: Language(LANG_LIB_PATH, lang_name) 
        for lang_name in set(EXTENSION_TO_LANGUAGE.values())
}

def get_language_for_extension(ext: str) -> Tuple(Path, str):
    lang_name = EXTENSION_TO_LANGUAGE.get(ext)
    return LANGUAGES.get(lang_name), lang_name
