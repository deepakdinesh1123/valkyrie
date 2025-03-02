from enum import Enum
from typing import List, Optional
from uuid import UUID

from pydantic import BaseModel, Field

from .base import ResponseBase


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
    env: Optional[dict[str, str]] = Field(
        default=None, description="Environment variables to be added"
    )


class NewTerminalResponse(ResponseBase):
    msgType: str = Field(default="NewTerminalResponse", frozen=True)
    terminalID: UUID = Field(description="Terminal ID")


class TerminalBaseModel(BaseModel):
    """Base model that stores terminal details"""

    terminalId: UUID = Field(
        description="Unique identifier for the terminal session", frozen=True
    )


class TerminalWrite(TerminalBaseModel):
    """Write to terminal"""

    msgType: str = Field(default="TerminalWrite", frozen=True)
    input: str = Field(..., description="input to write to terminal")


class TerminalWriteResponse(TerminalBaseModel, ResponseBase):
    """Response after writing to terminal"""

    msgType: str = Field(default="TerminalWriteResponse", frozen=True)


class TerminalRead(TerminalBaseModel):
    """Read from terminal"""

    msgType: str = Field(default="TerminalRead", frozen=True)


class TerminalReadResponse(TerminalBaseModel, ResponseBase):
    """Response after reading from terminal"""

    msgType: str = Field(default="TerminalReadResponse", frozen=True)
    output: str = Field(..., description="Content read from the terminal")
    eof: bool = Field(
        default=False, description="Indicates if the end of the stream has been reached"
    )


class TerminalClose(TerminalBaseModel):
    """Close terminal session"""

    msgType: str = Field(default="TerminalClose", frozen=True)


class TerminalCloseResponse(TerminalBaseModel, ResponseBase):
    """Response after closing a terminal session"""

    msgType: str = Field(default="TerminalCloseResponse", frozen=True)
