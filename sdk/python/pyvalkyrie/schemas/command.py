from typing import List
from pydantic import BaseModel, Field
from .terminal import EnvironmentVariable

class ExecuteCommand(BaseModel):
    """Execute a command in the sandbox environment"""
    msgType: str = Field(default="ExecuteCommand", frozen=True)
    sandboxId: int = Field(..., description="ID of the sandbox to execute the command in")
    command: str = Field(..., description="Command to execute")
    stdin: bool = Field(default=False, description="Enable stdin")
    stderr: bool = Field(default=True, description="Enable stderr")
    stdout: bool = Field(default=True, description="Enable stdout")
    workDir: str = Field(default="", description="Working directory for command execution")
    env: List[EnvironmentVariable] = Field(default_factory=list, description="Environment variables")