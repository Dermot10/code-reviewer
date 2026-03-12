import pytest
import json
from unittest.mock import AsyncMock, patch
from backend_python.worker.tasks import BATCH_SIZE, handle_review_task, handle_enhance_task, handle_assistant_task


@pytest.mark.asyncio
async def test_review_handler_sets_result_and_publishes():

    mock_redis = AsyncMock()

    task = type("Task", (), {"review_id": 1})()

    with patch(
        "backend_python.worker.tasks.process_review_task",
        AsyncMock(return_value={"score": 9})
    ):
        await handle_review_task(mock_redis, task)

    mock_redis.set.assert_called_once()
    mock_redis.publish.assert_called_once()

    #first argument set is the key 
    key = mock_redis.set.call_args[0][0]
    assert key == "review:1:result"



@pytest.mark.asyncio
async def test_enhance_handler_sets_result_and_publishes():

    mock_redis = AsyncMock()

    task = type("Task", (), {"enhancement_id": 2})()

    with patch(
        "backend_python.worker.tasks.process_enhance_task",
        AsyncMock(return_value={"improvement": "refactor"})
    ):
        await handle_enhance_task(mock_redis, task)

    mock_redis.set.assert_called_once()
    mock_redis.publish.assert_called_once()

    key = mock_redis.set.call_args[0][0]
    assert key == "enhancement:2:result"



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
async def test_assistant_batches_chunks_precise():

    mock_redis = AsyncMock()

    task = type("Task", (), {"conversation_id": 1, "user_id": 1})()
    
    # 1-char chunks yielded to sim stremaing 
    async def fake_stream(_):
        for _ in range(10):
            yield "x"  

    with patch(
        "backend_python.worker.tasks.process_assistant_task",
        fake_stream
    ):
        await handle_assistant_task(mock_redis, task)

    # gather all published payloads
    published_payloads = [
        json.loads(call.args[1])
        for call in mock_redis.publish.call_args_list
    ]

    # Check batching logic, each payload has 'chunk' <= batch_size or final done=True
    for payload in published_payloads[:-1]:  # skip final done=True
        chunk = payload["payload"]["chunk"]
        assert len(chunk) <= BATCH_SIZE

    # Final done=True event, last payload contains full result which is published to redis (mock)
    final_payload = published_payloads[-1]["payload"]
    assert final_payload["done"] is True
    assert final_payload["chunk"] == "x" * 10