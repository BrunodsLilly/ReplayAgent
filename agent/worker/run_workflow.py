import asyncio
from temporalio.client import Client

from workflows import PromptWorkflow


async def main():
    # Create client connected to server at the given address
    client = await Client.connect("localhost:7233")

    # Execute a workflow
    question = "How old is London?"
    result = await client.execute_workflow(
        PromptWorkflow.run, question, id="my-workflow-id", task_queue="my-task-queue"
    )

    print(f"Result: {result}")


if __name__ == "__main__":
    asyncio.run(main())
