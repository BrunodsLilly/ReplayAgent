import asyncio
from temporalio.client import Client
from temporalio.worker import Worker

# Import the activity and workflow from our other files
from activities import get_plan
from workflows import PromptWorkflow


async def main():
    # Create client connected to server at the given address
    client = await Client.connect("localhost:7233")

    # Run the worker
    worker = Worker(
        client,
        task_queue="my-task-queue",
        workflows=[PromptWorkflow],
        activities=[get_plan],
    )
    await worker.run()


if __name__ == "__main__":
    asyncio.run(main())
