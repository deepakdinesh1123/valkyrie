# /// script
# requires-python = ">=3.12"
# dependencies = [
#     "pydantic",
# ]
# ///

import json
from enum import Enum
from typing import List, Literal, Optional
from pathlib import Path

from pydantic import BaseModel, Field, ConfigDict

SCHEMAS_DIR = "schemas/agent"

# Constants
class Shell(str, Enum):
    NIX = "nix"
    BASH = "bash"
    SH = "sh"
    NIX_SHELL = "nix-shell"

class EnvironmentVariable(BaseModel):
    """Environment variable configuration"""
    key: str = Field(..., description="Environment variable name")
    value: str = Field(..., description="Environment variable value")

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

class NewTerminal(BaseModel):
    """Create a new terminal session"""
    msgType: str = Field(default="NewTerminal", frozen=True)
    shell: Shell = Field(..., description="Shell type to use")
    nix_shell: str = Field(default="", description="Nix shell configuration")
    nix_flake: str = Field(default="", description="Nix flake configuration")
    packages: List[str] = Field(default_factory=list, description="Packages to install")

class TerminalWrite(BaseModel):
    """Write to terminal"""
    msgType: str = Field(default="TerminalWrite", frozen=True)
    content: str = Field(..., description="Content to write to terminal")

class TerminalRead(BaseModel):
    """Read from terminal"""
    msgType: str = Field(default="TerminalRead", frozen=True)
    timeout: Optional[float] = Field(default=None, description="Read timeout in seconds")

class TerminalClose(BaseModel):
    """Close terminal session"""
    msgType: str = Field(default="TerminalClose", frozen=True)

class AddFile(BaseModel):
    """Add a file to the sandbox"""
    msgType: str = Field(default="AddFile", frozen=True)
    sandboxId: int = Field(..., description="ID of the sandbox to add the file to")
    file_name: str = Field(..., description="Name of the file")
    path: str = Field(..., description="Path where to create the file")
    content: str = Field(..., description="File content")

def ensure_directory(directory: str | Path) -> Path:
    """Ensure the directory exists and return Path object"""
    path = Path(directory)
    path.mkdir(parents=True, exist_ok=True)
    return path

def write_schema(directory: str | Path, model: type[BaseModel]) -> None:
    """Write schema for a model to a file"""
    dir_path = ensure_directory(directory)
    file_name = f"{model.__name__.lower()}.json"
    schema_path = dir_path / file_name
    
    schema = model.model_json_schema()
    
    with schema_path.open('w') as f:
        json.dump(schema, f, indent=4)
        f.write('\n')  # Add newline at end of file

def generate():
    """Generate all schemas"""
    schemas_dir = Path(SCHEMAS_DIR)
    
    models = [
        ExecuteCommand,
        AddFile,
        NewTerminal,
        TerminalWrite,
        TerminalRead,
        TerminalClose
    ]
    
    for model in models:
        write_schema(schemas_dir, model)
        print(f"Generated schema for {model.__name__}")

if __name__ == "__main__":
    generate()