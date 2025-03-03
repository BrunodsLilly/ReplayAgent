package replayagent

type ReplayAgentInput struct {
	Name string
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
	Prompt string
}
