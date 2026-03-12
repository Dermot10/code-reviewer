import pytest
from unittest.mock import AsyncMock, Mock
from backend_python.service.agent_service import agent_service, aggregate_reviews
from backend_python.schemas.ai.review_context import ResponseContext

@pytest.mark.asyncio
async def test_agent_service_calls_agents_and_aggregates():
    # Mock agents
    async def agent1(context):
        return [Mock(spec=ResponseContext, syntax=Mock(feedback="syntax1", issues=[]), semantics=None, security=None)]
    
    async def agent2(context):
        return [Mock(spec=ResponseContext, syntax=Mock(feedback="syntax2", issues=[]), semantics=None, security=None)]

    aggregate_func = Mock(return_value={"feedback": "combined", "issues": []})

    result = await agent_service([agent1, agent2], [{"code": "print('hi')"}], aggregate_func)

    # Each agent was called
    # aggregate_func called with combined results
    assert aggregate_func.call_count == 1
    assert result == {"feedback": "combined", "issues": []}


@pytest.mark.asyncio
async def test_agent_service_handles_partial_failure():
    async def agent1(context):
        raise ValueError("fail")
    
    async def agent2(context):
        return [Mock(spec=ResponseContext, syntax=Mock(feedback="ok", issues=[]), semantics=None, security=None)]

    aggregate_func = Mock(return_value={"feedback": "ok", "issues": []})

    result = await agent_service([agent1, agent2], [{"code": "print('hi')"}], aggregate_func)
    
    # Still returns aggregation from successful agent
    assert result == {"feedback": "ok", "issues": []}