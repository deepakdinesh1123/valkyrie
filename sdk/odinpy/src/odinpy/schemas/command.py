from typing import List, Literal

from pydantic import BaseModel, Field

from .base import ResponseBase


class Command(BaseModel):
    commandId: str = Field(..., description="Command ID")


class ExecuteCommand(BaseModel):
    """Execute a command in the sandbox environment"""

    command: str = Field(description="Command")
    msgType: str = Field(default="ExecuteCommand", frozen=True)
    stdin: bool = Field(default=False, description="Enable stdin")
    stderr: bool = Field(default=True, description="Enable stderr")
    stdout: bool = Field(default=True, description="Enable stdout")
    workDir: str = Field(
        default="", description="Working directory for command execution"
    )
    env: List[dict[str, str]] = Field(
        default_factory=list, description="Environment variables"
    )


class ExecuteCommandResponse(Command, ResponseBase):
    msgType: str = Field(default="ExecuteCommandResponse", frozen=True)
    state: Literal["running", "exited", "stopped", "starting"] = Field(
        default="starting"
    )
    stdout: str = Field(..., description="stdout")


class CommandWriteInput(Command):
    msgType: str = Field(default="CommandWriteInput", frozen=True)
    input: str = Field(default="input")


class CommandWriteInputResponse(Command, ResponseBase):
    msgType: str = Field(default="CommandWriteInputResponse", frozen=True)


class CommandReadOutput(Command):
    msgType: str = Field(default="CommandReadOutput", frozen=True)


class CommandReadOutputResponse(Command, ResponseBase):
    msgType: str = Field(default="CommandReadOutputResponse", frozen=True)
    stdout: str


class CommandTerminate(Command):
    msgType: str = Field(default="CommandTerminate", frozen=True)


class CommandTerminateResponse(Command, ResponseBase):
    msgType: str = Field(default="CommandTerminateResponse", frozen=True)
