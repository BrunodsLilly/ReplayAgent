package replayagent

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

var activities *Activities

// Workflow is a ReplayAgent workflow definition.
func ReplayAgentWF(ctx workflow.Context, input ReplayAgentInput) (ReplayAgentOutput, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("ReplayAgent workflow started", "input", input)

	var result HelloActivityOutput
	err := workflow.ExecuteActivity(ctx, activities.HelloActivity, HelloActivityInput{Name: input.Name}).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return ReplayAgentOutput{}, err
	}

	var promptLLMResult PromptLLMActivityOutput
	question := "Recommend a list of a list of tools and parameters I need to garden my backyard."
	err = workflow.ExecuteActivity(ctx, activities.PromptLLMActivity, PromptLLMActivityInput{Prompt: question}).Get(ctx, &promptLLMResult)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return ReplayAgentOutput{}, err
	}

	logger.Info("PromptLLMActivity", "result", promptLLMResult)

	logger.Info("ReplayAgent workflow completed.", "result", result)

	return ReplayAgentOutput{Result: result.Result}, nil
}

func WorkflowID(id string) string {
	return "replay_agent_" + id
}
