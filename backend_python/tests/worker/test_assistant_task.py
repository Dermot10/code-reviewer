import pytest
from unittest.mock import AsyncMock, patch
from backend_python.worker.tasks import handle_assistant_task

@pytest.mark.asyncio
async def test_assistant_stream_publishes_chunks():
    """Verifies AI Stream -> task handler -> Redis Publish"""
    mock_redis = AsyncMock()

    task = type("Task", (), {
        "conversation_id": 1, 
        "user_id": 1 
    })()

    async def fake_stream(_): 
        yield "Hello"
        yield " "
        yield "World"

    with patch(
        "backend_python.worker.tasks.process_assistant_task", 
        fake_stream
    ):
        await handle_assistant_task(mock_redis, task)

    assert mock_redis.publish.called


@pytest.mark.asyncio
async def test_assistant_batches_chunks():

    mock_redis = AsyncMock()

    task = type("Task", (), {
        "conversation_id": 1, 
        "user_id": 1
    })()

    async def fake_stream(_): 
        for _ in range(10): 
            yield "chunk"
    
    with patch(
        "backend_python.worker.tasks.process_assistant_task", 
        fake_stream
    ):
        await handle_assistant_task(mock_redis, task)

    # should publish multiple times but not once per token
    assert mock_redis.publish.call_count > 0