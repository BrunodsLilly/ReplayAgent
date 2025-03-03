from pydantic import BaseModel

class ToolParams(BaseModel):
    tool: str
    params: dict

class PromptResponse(BaseModel):
    tool_params: list[list[ToolParams]]
