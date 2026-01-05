from enum import Enum



# migrate to go service
class ExportType(Enum): 
    MD = "markdown"
    TXT = "txt"
    CSV = "csv"
    JSON = "json"
    PY = "py"