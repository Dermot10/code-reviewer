import pytest
from unittest.mock import AsyncMock, patch
from backend_python.service.review_service import Execute_review

@pytest.mark.asyncio
async def test_execute_review_returns_aggregated_output():
    chunked_code = [{"code": "print('hello')"}]  # dummy CodeContext-like dict

    # Mock the underlying code_review_service
    fake_result = {"feedback": "ok", "issues": []}
    with patch("backend_python.service.review_service.code_review_service", AsyncMock(return_value=fake_result)):
        result = await Execute_review(chunked_code)

    assert result == fake_result


@pytest.mark.asyncio
async def test_execute_review_logs_exception():
    chunked_code = [{"code": "bad_code"}]

    with patch("backend_python.service.review_service.code_review_service", AsyncMock(side_effect=ValueError("fail"))), \
         patch("backend_python.service.review_service.logger") as mock_logger:
        with pytest.raises(ValueError):
            await Execute_review(chunked_code)
        mock_logger.warning.assert_called()