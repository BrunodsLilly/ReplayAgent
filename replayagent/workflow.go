package replayagent

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

// Workflow is a ReplayAgent workflow definition.
func ReplayAgentWF(ctx workflow.Context, input ReplayAgentInput) (ReplayAgentOutput, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("ReplayAgent workflow started", "input", input)

	var result HelloActivityOutput
	err := workflow.ExecuteActivity(ctx, HelloActivity, HelloActivityInput{Name: input.Name}).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return ReplayAgentOutput{}, err
	}

	logger.Info("ReplayAgent workflow completed.", "result", result)

	return ReplayAgentOutput{Result: result.Result}, nil
}

func HelloActivity(ctx context.Context, input HelloActivityInput) (HelloActivityOutput, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "name", input.Name)
	return HelloActivityOutput{Result: "ReplayAgent " + input.Name + "!"}, nil
}

func WorkflowID(id string) string {
	return "replay_agent_" + id
}
