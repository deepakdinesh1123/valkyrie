from enum import Enum
from typing import List, Optional
from uuid import UUID

from pydantic import BaseModel, Field


class EnvironmentVariable(BaseModel):
    """Environment variable configuration"""

    key: str = Field(..., description="Environment variable name")
    value: str = Field(..., description="Environment variable value")


class Shell(str, Enum):
    NIX = "nix"
    BASH = "bash"
    SH = "sh"
    NIX_SHELL = "nix-shell"


class NewTerminal(BaseModel):
    """Create a new terminal session"""

    msgType: str = Field(default="NewTerminal", frozen=True)
    shell: Optional[Shell] = Field(default=None, description="Shell type to use")
    nixShell: Optional[str] = Field(default=None, description="Nix shell configuration")
    nixFlake: Optional[str] = Field(default=None, description="Nix flake configuration")
    packages: Optional[List[str]] = Field(
        default=None, description="Packages to install"
    )
    env: Optional[List[EnvironmentVariable]] = Field(
        default=None, description="Environment variables to be added"
    )


class NewTerminalResponse(BaseModel):
    terminalID: UUID = Field(description="Terminal ID")
    msg: str = Field(..., description="Message")
    success: bool = Field(description="Success")


class TerminalBaseModel(BaseModel):
    """Base model that stores terminal details"""

    terminalId: UUID = Field(
        description="Unique identifier for the terminal session", frozen=True
    )


class TerminalWrite(TerminalBaseModel):
    """Write to terminal"""

    msgType: str = Field(default="TerminalWrite", frozen=True)
    input: str = Field(..., description="input to write to terminal")


class TerminalWriteResponse(TerminalBaseModel):
    """Response after writing to terminal"""

    msg: str = Field(..., description="Message confirming write operation")
    success: bool = Field(description="Success")


class TerminalRead(TerminalBaseModel):
    """Read from terminal"""

    msgType: str = Field(default="TerminalRead", frozen=True)


class TerminalReadResponse(TerminalBaseModel):
    """Response after reading from terminal"""

    output: str = Field(..., description="Content read from the terminal")
    eof: bool = Field(
        default=False, description="Indicates if the end of the stream has been reached"
    )
    success: bool = Field(description="Success")
    msg: str = Field(..., description="optional message")


class TerminalClose(TerminalBaseModel):
    """Close terminal session"""

    msgType: str = Field(default="TerminalClose", frozen=True)


class TerminalCloseResponse(TerminalBaseModel):
    """Response after closing a terminal session"""

    msg: str = Field(..., description="Message confirming terminal closure")
    success: bool = Field(description="Success")
