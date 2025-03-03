from pydantic_ai import Agent
from temporalio import activity

from shared import PromptResponse, Service


service_registry = [
    Service(
        name="set_temperature",
        address="http://localhost:8080",
        params={"temperature": "$temperature"},
    ),
    Service(
        name="set_pressure",
        address="http://localhost:8080",
        params={"pressure": "$pressure"},
    ),
    Service(
        name="set_flow_rate",
        address="http://localhost:8080",
        params={"flow_rate": "$flow_rate"},
    ),
]


@activity.defn
async def get_plan(prompt: str) -> PromptResponse:
    agent = Agent("openai:gpt-4o", result_type=PromptResponse)
    service_context = """Available services:
        - set_temperature: {$temperature}
        - set_pressure: {$pressure}
        - set_flow_rate: {$flow_rate}

        Replace the placeholders with the desired values.
    """
    prompt = f"{prompt}\n{service_context}"

    result_sync = await agent.run(prompt)
    return result_sync.data
