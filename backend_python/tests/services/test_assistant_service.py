import pytest
from unittest.mock import AsyncMock, patch
from backend_python.service.assistant_service import Execute_assistant

@pytest.mark.asyncio
async def test_execute_assistant_yields_chunks():
    # Mock handle_assistant to yield chunks
    async def fake_assistant(prompt):
        yield "Hello"
        yield "World"

    with patch("backend_python.service.assistant_service.handle_assistant", new=fake_assistant):
        chunks = []
        async for chunk in Execute_assistant("Hi"):
            chunks.append(chunk)

    assert chunks == ["Hello", "World"]

@pytest.mark.asyncio
async def test_execute_assistant_logs_exception():
    """Ensure exceptions in handle_assistant are logged and propagated."""

    # Async generator that immediately raises
    async def failing_assistant(prompt):
        raise ValueError("fail")
        yield  # Needed to make it an async generator (won't actually run)

    with patch("backend_python.service.assistant_service.handle_assistant", new=failing_assistant), \
         patch("backend_python.service.assistant_service.logger") as mock_logger:
        
        with pytest.raises(ValueError):
            async for _ in Execute_assistant("Hi"):
                pass

        # Logger should be called
        mock_logger.warning.assert_called()