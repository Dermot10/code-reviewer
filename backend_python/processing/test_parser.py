# test_parser.py

from typing import List
from backend_python.processing.preprocessing import extract_chunks
from backend_python.processing.context import CodeContext 

# Path to the sample file
file_path = "backend_python/processing/sample.py"

# Read the file
with open(file_path, "r", encoding="utf-8") as f:
    code = f.read()

# Extract chunks
chunks: List[CodeContext] = extract_chunks(file_path, code)

# Print results
for chunk in chunks:
    print(f"CHUNK ID: {chunk.chunk_id}")
    print(f"LINES:\n{chunk.code}")
    print("-" * 40)
