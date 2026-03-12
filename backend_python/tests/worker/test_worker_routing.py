import pytest
from unittest.mock import AsyncMock, patch
from backend_python.worker.worker import handle_task

@pytest.fixture
def assistant_task():
    return {
        "type": "assistant",
        "conversation_id": 1,
        "user_id": 1,
        "prompt": "Hello World"
    }

@pytest.fixture
def review_task():
    return {
        "type": "review",
        "review_id": 2,
        "user_id": 42,
        "code": "print('hi')"
    }

@pytest.fixture
def enhance_task():
    return {
        "type": "enhance",
        "enhancement_id": 3,
        "user_id": 99,
        "code": "print('hi')"
    }



@pytest.mark.asyncio
async def test_worker_routes_assistant_task(assistant_task):

    with patch.dict("os.environ", {"LOCK_EXPIRY": "30"}), \
         patch("backend_python.worker.worker.r", new_callable=AsyncMock) as mock_redis, \
         patch("backend_python.worker.worker.handle_assistant_task", new_callable=AsyncMock) as mock_assistant:

        mock_redis.set.return_value = True
        mock_redis.delete.return_value = True

        await handle_task(assistant_task)

        mock_assistant.assert_called_once()

        mock_redis.set.assert_called_with(
            "assistant:1:lock", "1", nx=True, ex=30
        )


@pytest.mark.asyncio
async def test_worker_routes_review_task(review_task):

    with patch.dict("os.environ", {"LOCK_EXPIRY": "30"}), \
         patch("backend_python.worker.worker.r", new_callable=AsyncMock) as mock_redis, \
         patch("backend_python.worker.worker.handle_review_task", new_callable=AsyncMock) as mock_review:

        mock_redis.set.return_value = True
        mock_redis.delete.return_value = True

        await handle_task(review_task)

        mock_review.assert_called_once()

        mock_redis.set.assert_called_with(
            "review:2:lock", "1", nx=True, ex=30
        )



@pytest.mark.asyncio
async def test_worker_skips_if_lock_exists(assistant_task):

    with patch("backend_python.worker.worker.r", new_callable=AsyncMock) as mock_redis:

        mock_redis.set.return_value = False

        await handle_task(assistant_task)

        mock_redis.delete.assert_not_called()



@pytest.mark.asyncio
async def test_worker_publishes_failure_event(assistant_task):

    with patch("backend_python.worker.worker.r", new_callable=AsyncMock) as mock_redis, \
         patch(
             "backend_python.worker.worker.handle_assistant_task",
             new_callable=AsyncMock,
             side_effect=Exception("boom")
         ):

        mock_redis.set.return_value = True
        mock_redis.delete.return_value = True

        await handle_task(assistant_task)

        mock_redis.publish.assert_called()