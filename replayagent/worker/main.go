package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"app/replayagent"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	agentTemporalClient, err := client.Dial(client.Options{
		HostPort: "localhost:7233",
	})
	if err != nil {
		log.Fatalln("Unable to create agent client", err)
	}
	defer agentTemporalClient.Close()

	w := worker.New(c, "replay-agent", worker.Options{})

	activities := replayagent.NewActivities(agentTemporalClient)

	w.RegisterWorkflow(replayagent.ReplayAgentWF)
	w.RegisterActivity(activities)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
