from enum import Enum
from uuid import UUID
from typing import List, Optional
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
    shell: Shell = Field(..., description="Shell type to use")
    nixShell: str = Field(default="", description="Nix shell configuration")
    nixFlake: str = Field(default="", description="Nix flake configuration")
    packages: List[str] = Field(default_factory=list, description="Packages to install")
    env: List[EnvironmentVariable] = Field(default_factory=list, description="Environment variables to be added")

class NewTerminalResponse(BaseModel):
    terminalID: UUID = Field(description="Terminal ID")
    msg: str = Field(..., description="Message")
    success: bool = Field(description="Success")

class TerminalWrite(BaseModel):
    """Write to terminal"""
    msgType: str = Field(default="TerminalWrite", frozen=True)
    terminalId: UUID = Field(..., description="Unique identifier for the terminal session")
    input: str = Field(..., description="input to write to terminal")

class TerminalWriteResponse(BaseModel):
    """Response after writing to terminal"""
    terminalId: UUID = Field(..., description="Unique identifier for the terminal session")
    msg: str = Field(..., description="Message confirming write operation")
    success: bool = Field(description="Success")

class TerminalRead(BaseModel):
    """Read from terminal"""
    msgType: str = Field(default="TerminalRead", frozen=True)
    terminalId: UUID = Field(..., description="Unique identifier for the terminal session")
    timeout: Optional[float] = Field(default=None, description="Read timeout in seconds")

class TerminalReadResponse(BaseModel):
    """Response after reading from terminal"""
    terminalId: UUID = Field(..., description="Unique identifier for the terminal session")
    output: str = Field(..., description="Content read from the terminal")
    eof: bool = Field(default=False, description="Indicates if the end of the stream has been reached")
    success: bool = Field(description="Success")
    msg: str = Field(..., description="optional message")

class TerminalClose(BaseModel):
    """Close terminal session"""
    msgType: str = Field(default="TerminalClose", frozen=True)
    terminalId: UUID = Field(..., description="Unique identifier for the terminal session")

class TerminalCloseResponse(BaseModel):
    """Response after closing a terminal session"""
    terminalId: UUID = Field(..., description="Unique identifier for the terminal session")
    msg: str = Field(..., description="Message confirming terminal closure")
    success: bool = Field(description="Success")