from pydantic_ai import Agent
from temporalio import activity


@activity.defn
async def ask_question(prompt: str) -> str:
    agent = Agent("openai:gpt-4o")

    result_sync = await agent.run(prompt)
    return result_sync.data
