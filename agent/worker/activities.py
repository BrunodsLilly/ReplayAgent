from pydantic_ai import Agent
from temporalio import activity

from shared import PromptResponse


@activity.defn
async def get_plan(prompt: str) -> PromptResponse:
    agent = Agent("openai:gpt-4o", result_type=PromptResponse)
    result_sync = await agent.run(prompt)
    return result_sync.data
