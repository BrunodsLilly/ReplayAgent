package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"go.temporal.io/sdk/client"

	"app/replayagent"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// new uuid
	uuid := uuid.New().String()

	workflowOptions := client.StartWorkflowOptions{
		ID:        replayagent.WorkflowID(uuid),
		TaskQueue: "replay-agent",
	}

	we, err := c.ExecuteWorkflow(context.Background(),
		workflowOptions, replayagent.ReplayAgentWF,
		replayagent.ReplayAgentInput{Name: "Temporal"})
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// Synchronously wait for the workflow completion.
	var result replayagent.ReplayAgentOutput
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)
}
