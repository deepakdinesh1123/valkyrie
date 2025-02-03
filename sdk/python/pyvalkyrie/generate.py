# /// script
# requires-python = ">=3.12"
# dependencies = [
#     "pydantic",
# ]
# ///

import json

from pathlib import Path

from pydantic import BaseModel

from schemas.terminal import (
    NewTerminal,
    NewTerminalResponse,
    TerminalWrite,
    TerminalWriteResponse,
    TerminalRead,
    TerminalReadResponse,
    TerminalClose,
    TerminalCloseResponse
)

from schemas.command import (
    ExecuteCommand
)

from schemas.file import (
    AddFile
)

SCHEMAS_DIR = "schemas/agent"


def ensure_directory(directory: str | Path) -> Path:
    """Ensure the directory exists and return Path object"""
    path = Path(directory)
    path.mkdir(parents=True, exist_ok=True)
    return path

def write_schema(directory: str | Path, model: type[BaseModel]) -> None:
    """Write schema for a model to a file"""
    dir_path = ensure_directory(directory)
    file_name = f"{model.__name__}.json"
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
        NewTerminalResponse,
        TerminalWrite,
        TerminalWriteResponse,
        TerminalRead,
        TerminalReadResponse,
        TerminalClose,
        TerminalCloseResponse,
    ]
    
    for model in models:
        write_schema(schemas_dir, model)
        print(f"Generated schema for {model.__name__}")

if __name__ == "__main__":
    generate()