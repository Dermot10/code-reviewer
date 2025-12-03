# Workflow Notes

Request -> Handler -> Preprocessing
-> syntax agent
-> semantics agent
-> security agent
-> best practices agent
-> Aggregator func -> final agent -> Postprocessing -> Handler -> Response

# Tree Parsing Code

The file is the tree root, every top level class or func is a child of the root, and the methods inside are nested children

root (file)
├─ ClassDef: Greeter
│ ├─ FunctionDef: \__init_
│ └─ FunctionDef: greet
├─ FunctionDef: add
└─ AsyncFunctionDef: async_task

The AST library is used for parsing the python code. Like all efficient parsers dfs is used internally

BFS logic for ast walk-
done to identify all major chunks first,
then extract their code

## Algorithm

DFS algorithm looks to traverse the tree prioritsing a vertical traversal pattern.
With the next node being the child before moving to siblings.
Visiting the node first, then it's children

# Other-

pip install fastapi uvicorn openai pydantic dotenv prometheus-client python-multipart
