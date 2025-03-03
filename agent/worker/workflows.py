from datetime import timedelta
from temporalio import workflow

from shared import PromptResponse


# Import our activity, passing it through the sandbox
with workflow.unsafe.imports_passed_through():
    from activities import get_plan


@workflow.defn
class PromptWorkflow:
    """A workflow for asking a question to an LLM."""

    @workflow.run
    async def run(self, prompt: str) -> PromptResponse:
        return await workflow.execute_activity(
            get_plan, prompt, schedule_to_close_timeout=timedelta(seconds=50)
        )
