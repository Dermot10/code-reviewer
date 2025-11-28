from prometheus_client import Counter, Histogram

CHUNKING_FAILURES = Counter(
    "file_processing_failures_total", "Code input which has failed the chunking process"
)
FILE_EXTRACTION_ERRORS = Counter(
    "file_extraction_errors_total", "Number of files which have failed the chunking process"
) 

SYNTAX_ERRORS = Counter(
    "syntax_errors_total", "Number of OpenAI processing errors - when processing syntax"
)

SEMANTICS_ERRORS = Counter(
    "semantics_errors_total", "Number of OpenAI processing errors - when processing semantics"
)

BEST_PRACTICES_ERRORS = Counter(
    "best_practices_errors_total", "Number of OpenAI processing errors - when processing best practices"
)

SECURITY_ERRORS = Counter(
    "security_errors_total", "Number of OpenAI processing errors - when processing for security "
)

AGGREGATOR_ERRORS = Counter(
    "openai_errors_total", "Number of OpenAI processing errors"
)


AI_PROCESSING_TIME = Histogram(
    "ai_processing_seconds", "Time spent calling OpenAI API"
)
