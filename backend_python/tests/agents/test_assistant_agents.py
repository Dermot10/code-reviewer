import pytest
from unittest.mock import patch, MagicMock
from backend_python.ai.ai_assistant_agents import handle_assistant
from backend_python.exceptions.exceptions import OpenAiProcessingError
from backend_python.metrics import ASSISTANT_ERRORS



async def collect_async(gen):
    results = []
    async for item in gen:
        results.append(item)
    return results


@pytest.mark.asyncio
async def test_handle_assistant_streams_chunks():
    """Agent should stream chunks returned by the OpenAI client."""

    async def fake_stream(system_prompt, prompt):
        yield "hello"
        yield " world"

    with patch(
        "backend_python.ai.ai_assistant_agents.stream_open_ai_call",
        new=fake_stream,
    ):
        result = await collect_async(handle_assistant("hi"))

    assert result == ["hello", " world"]



@pytest.mark.asyncio
async def test_handle_assistant_raises_and_increments_metric():
    """Agent should increment metric and propagate OpenAiProcessingError."""

    async def failing_stream(system_prompt, prompt):
        # define an async generator that immediately raises
        if False:
            yield
        raise OpenAiProcessingError("fail")

    with patch(
        "backend_python.ai.ai_assistant_agents.stream_open_ai_call",
        new=failing_stream,
    ), patch.object(ASSISTANT_ERRORS, "inc") as mock_inc:

        with pytest.raises(OpenAiProcessingError):
            await collect_async(handle_assistant("hi"))

        mock_inc.assert_called_once()