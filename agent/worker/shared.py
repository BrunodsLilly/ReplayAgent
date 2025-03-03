from pydantic import BaseModel

class ToolParams(BaseModel):
    tool: str
    params: dict[str, str]
    
    # Add a method to convert to a dict format
    def to_dict(self):
        return {
            "tool": self.tool,
            "params": self.params
        }

class PromptResponse(BaseModel):
    tool_params: list[list[ToolParams]]
    
    # Add a method to convert to the format expected by Go
    def to_dict(self):
        return {
            "tool_params": [[tp.to_dict() for tp in group] for group in self.tool_params]
        }