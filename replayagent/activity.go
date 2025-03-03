package replayagent

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
)

type Activities struct {
	agentTemporalClient client.Client
}

func NewActivities(agentTemporalClient client.Client) *Activities {
	return &Activities{agentTemporalClient: agentTemporalClient}
}

func (a *Activities) HelloActivity(ctx context.Context, input HelloActivityInput) (HelloActivityOutput, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "name", input.Name)
	return HelloActivityOutput{Result: "ReplayAgent " + input.Name + "!"}, nil
}

func (a *Activities) PromptLLMActivity(ctx context.Context, input PromptLLMActivityInput) (PromptLLMActivityOutput, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("PromptLLMActivity", "input", input)

	workflowOptions := client.StartWorkflowOptions{
		ID:        "agent" + input.UUID.String(),
		TaskQueue: "my-task-queue",
	}

	we, err := a.agentTemporalClient.ExecuteWorkflow(ctx, workflowOptions, "PromptWorkflow", input.Prompt)
	if err != nil {
		logger.Error("Failed to execute workflow", "error", err)
		return PromptLLMActivityOutput{}, err
	}

	var result PromptLLMActivityOutput
	err = we.Get(ctx, &result)
	if err != nil {
		logger.Error("Failed to get workflow result", "error", err)
		return PromptLLMActivityOutput{}, err
	}

	logger.Info("PromptLLMActivity", "result", result)

	return result, nil
}

func (a *Activities) GetSensorReadings(ctx context.Context, input GetSensorReadingsInput) (GetSensorReadingsOutput, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("GetSensorReadings", "input", input)

	// based on the step, return the sensor readings
	// have a deterministic way of generating the sensor readings
	// generate the sensor readings based on the step
	var sensorReadings []SensorReading
	var types = []string{"temperature", "humidity", "light", "pressure"}
	for i := 0; i < 3; i++ {
		sensorReadings = append(sensorReadings, SensorReading{
			SensorType:  types[i],
			SensorValue: strconv.Itoa(input.Step*10 + i),
		})
	}

	return GetSensorReadingsOutput{
		SensorReadings: sensorReadings,
	}, nil
}

func (a *Activities) CheckSuccess(ctx context.Context, input CheckSuccessInput) (CheckSuccessOutput, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("CheckSuccess", "input", input)

	if input.Step == 0 {
		return CheckSuccessOutput{Success: false}, nil
	}

	success := (input.Step%3 == 0) || (input.Step > 7)

	logger.Info("CheckSuccess result", "step", input.Step, "success", success)
	return CheckSuccessOutput{Success: success}, nil
}

func (a *Activities) ExecuteTool(ctx context.Context, input ExecuteToolInput) (ExecuteToolOutput, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("ExecuteTool", "input", input)

	// based on the step, sleep for a random amount of time between 1 and 4 seconds
	time.Sleep(time.Duration(rand.Intn(3)+1) * time.Second)
	return ExecuteToolOutput{Success: true}, nil
}
