import pytest
from unittest.mock import patch
from backend_python.ai.ai_review_agents import (
    handle_syntax,
    handle_sematics,
    handle_best_practices,
    handle_security,
)
from backend_python.schemas.ai.review_context import CodeContext, ResponseContext
from backend_python.metrics import SYNTAX_ERRORS, SEMANTICS_ERRORS, BEST_PRACTICES_ERRORS, SECURITY_ERRORS
from backend_python.exceptions.exceptions import OpenAiProcessingError

@pytest.mark.asyncio
@pytest.mark.parametrize(
    "agent_func, strategy_field",
    [
        (handle_syntax, "syntax"),
        (handle_sematics, "semantics"),
        (handle_best_practices, "best_practices"),
        (handle_security, "security"),
    ],
)
async def test_review_agents_happy_path(agent_func, strategy_field):
    """Agents should return ResponseContext with the correct strategy field."""

    # Dummy code context
    code_contexts = [CodeContext(chunk_id="1", file_path="fake.py", code="print('hi')")]

    # Dummy response for the strategy field
    dummy_response = {"feedback": "ok"}

    # Patch handle_agent to return ResponseContext with the correct field
    async def fake_handle_agent(code_contexts, *args, **kwargs):
        results = []
        for ctx in code_contexts:
            rc = ResponseContext(chunk_id=ctx.chunk_id)
            setattr(rc, strategy_field, dummy_response)
            results.append(rc)
        return results

    with patch("backend_python.ai.ai_review_agents.handle_agent", new=fake_handle_agent):
        results = await agent_func(code_contexts)

    # Assertions
    assert isinstance(results, list)
    assert len(results) == len(code_contexts)
    for r in results:
        assert isinstance(r, ResponseContext)
        assert getattr(r, strategy_field) == dummy_response



@pytest.mark.asyncio
@pytest.mark.parametrize(
    "agent_func, metric",
    [
        (handle_syntax, SYNTAX_ERRORS),
        (handle_sematics, SEMANTICS_ERRORS),
        (handle_best_practices, BEST_PRACTICES_ERRORS),
        (handle_security, SECURITY_ERRORS),
    ],
)
async def test_review_agents_raises_and_increments_metric(agent_func, metric):
    """Agents should increment the correct metric and propagate OpenAiProcessingError."""

    code_contexts = [CodeContext(chunk_id="1", file_path="fake.py", code="print('hi')")]

    # Async function that immediately raises
    async def failing_handle_agent(*args, **kwargs):
        raise OpenAiProcessingError("fail")

    with patch("backend_python.ai.ai_review_agents.handle_agent", new=failing_handle_agent), \
         patch.object(metric, "inc") as mock_inc:

        with pytest.raises(OpenAiProcessingError):
            await agent_func(code_contexts)

        # Metric should be incremented once
        mock_inc.assert_called_once()