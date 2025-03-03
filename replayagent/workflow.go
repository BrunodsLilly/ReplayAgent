package replayagent

import (
	"time"

	"go.temporal.io/sdk/temporal"
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

	// take in a prompt
	// execute the prompt
	// get the tool params
	// execute the tool
	// check if the check success is true, if not, reprompt
	// run the next step
	// if check success is true, return the result

	var toolParams [][]ToolParams
	var step int

	var sensorReadingsOutput GetSensorReadingsOutput
	err := workflow.ExecuteActivity(ctx, activities.GetSensorReadings,
		GetSensorReadingsInput{Step: step}).Get(ctx, &sensorReadingsOutput)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return ReplayAgentOutput{}, err
	}

	var promptLLMResult PromptLLMActivityOutput
	question := "Recommend a list of a list of tools and parameters I need to garden my backyard."
	err = workflow.ExecuteActivity(ctx, activities.PromptLLMActivity,
		PromptLLMActivityInput{
			Prompt:         question,
			UUID:           input.UUID,
			ToolParams:     toolParams,
			SensorReadings: sensorReadingsOutput.SensorReadings,
		}).Get(ctx, &promptLLMResult)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return ReplayAgentOutput{}, err
	}

	var executeToolResult ExecuteToolOutput
	var checkSuccessResult CheckSuccessOutput

	for {
		toolParams = promptLLMResult.ToolParams
		childCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			StartToCloseTimeout:    time.Second * 30,
			ScheduleToCloseTimeout: time.Minute * 5,
			RetryPolicy: &temporal.RetryPolicy{
				InitialInterval: time.Millisecond * 100,
				MaximumAttempts: 5,
			},
			Summary: toolParams[step][0].Tool,
		})

		err = workflow.ExecuteActivity(childCtx, activities.ExecuteTool,
			ExecuteToolInput{
				Step: step,
				Tool: toolParams[step][0].Tool,
			}).Get(ctx, &executeToolResult)
		if err != nil {
			logger.Error("Activity failed.", "Error", err)
			return ReplayAgentOutput{}, err
		}

		err = workflow.ExecuteActivity(ctx, activities.CheckSuccess,
			CheckSuccessInput{
				Step: step,
			}).Get(ctx, &checkSuccessResult)
		if err != nil {
			logger.Error("Activity failed.", "Error", err)
			return ReplayAgentOutput{}, err
		}

		var sensorReadingsOutput GetSensorReadingsOutput
		err := workflow.ExecuteActivity(ctx, activities.GetSensorReadings,
			GetSensorReadingsInput{Step: step}).Get(ctx, &sensorReadingsOutput)
		if err != nil {
			logger.Error("Activity failed.", "Error", err)
			return ReplayAgentOutput{}, err
		}

		logger.Info("CheckSuccess result", "step", step, "success", checkSuccessResult.Success)

		if checkSuccessResult.Success || step >= len(toolParams) {
			break
		}

		step++
	}

	return ReplayAgentOutput{Result: "success"}, nil
}

func WorkflowID(id string) string {
	return "replay_agent_" + id
}
