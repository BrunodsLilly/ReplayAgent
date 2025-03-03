from datetime import timedelta
from temporalio import workflow


# Import our activity, passing it through the sandbox
with workflow.unsafe.imports_passed_through():
    from activities import ask_question


@workflow.defn
class PromptWorkflow:
    """A workflow for asking a question to an LLM."""

    @workflow.run
    async def run(self, prompt: str) -> str:
        return await workflow.execute_activity(
            ask_question, prompt, schedule_to_close_timeout=timedelta(seconds=5)
        )
