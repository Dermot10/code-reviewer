# backend_python/tests/worker/test_worker_routing.py
import pytest
import json
from unittest.mock import AsyncMock, patch
from backend_python.worker.worker import handle_task

@pytest.mark.asyncio
async def test_worker_task_routing():
    # sample tasks for all types
    tasks = [
        {
            "type": "assistant",
            "conversation_id": 1,
            "user_id": 1,
            "prompt": "Hello World"
        },
        {
            "type": "review",
            "review_id": 2,
            "user_id": 42,
            "code": "print('hi')"
        },
        {
            "type": "enhance",
            "enhancement_id": 3,
            "user_id": 99,
            "code": "print('hi')"
        }
    ]

    # patch Redis client and task handlers inside worker.py
    with patch.dict("os.environ", {"LOCK_EXPIRY": "30"}) , \
        patch("backend_python.worker.worker.r", new_callable=AsyncMock) as mock_redis, \
        patch("backend_python.worker.worker.handle_assistant_task", new_callable=AsyncMock) as mock_assistant, \
        patch("backend_python.worker.worker.handle_review_task", new_callable=AsyncMock) as mock_review, \
        patch("backend_python.worker.worker.handle_enhance_task", new_callable=AsyncMock) as mock_enhance:

        # Redis mocks
        mock_redis.set.return_value = True
        mock_redis.delete.return_value = True
        mock_redis.publish.return_value = 1

        # Run tasks
        await handle_task(tasks[0])
        await handle_task(tasks[1])
        await handle_task(tasks[2])

        # Assert the correct handler was called once for each task
        mock_assistant.assert_called_once()
        mock_review.assert_called_once()
        mock_enhance.assert_called_once()

        # Check Redis locks called correctly
        mock_redis.set.assert_any_call("assistant:1:lock", "1", nx=True, ex=30)
        mock_redis.set.assert_any_call("review:2:lock", "1", nx=True, ex=30)
        mock_redis.set.assert_any_call("enhance:3:lock", "1", nx=True, ex=30)