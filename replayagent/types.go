package replayagent

import "github.com/google/uuid"

type ReplayAgentInput struct {
	Name string
	UUID uuid.UUID
}

type ReplayAgentOutput struct {
	Result string
}

type HelloActivityInput struct {
	Name string
}

type HelloActivityOutput struct {
	Result string
}

type PromptLLMActivityOutput struct {
	ToolParams [][]ToolParams `json:"tool_params"`
}

type ToolParams struct {
	Tool   string            `json:"tool"`
	Params map[string]string `json:"params"`
}

type PromptLLMActivityInput struct {
	UUID           uuid.UUID
	Prompt         string
	SensorReadings []SensorReading
	ToolParams     [][]ToolParams
}

type SensorReading struct {
	SensorType  string `json:"sensor_type"`
	SensorValue string `json:"sensor_value"`
}

type GetSensorReadingsInput struct {
	Step int `json:"step"`
}

type GetSensorReadingsOutput struct {
	SensorReadings []SensorReading `json:"sensor_readings"`
}

type CheckSuccessInput struct {
	Step int
}

type CheckSuccessOutput struct {
	Success bool
}

type ExecuteToolInput struct {
	Step int
	Tool string
}

type ExecuteToolOutput struct {
	Success bool
}
