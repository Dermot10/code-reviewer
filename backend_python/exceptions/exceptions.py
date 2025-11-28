class FileProcessingError(Exception):
    """Raised when a file is failed to be processed. 
    
    This includes; file validation, content validation and upload confirmation
    """
    def __init__(self, message: str): 
        super().__init__(message)

class FileExtractionError(Exception): 
    """Raised when an extraction method has failed to pull local file data"""
    def __init__(self, message: str): 
        super().__init__(message)


class OpenAiProcessingError(Exception): 
    """Raised when the openai model has failed to execute the call made"""
    def __init__(self, message: str): 
        super().__init__(message)