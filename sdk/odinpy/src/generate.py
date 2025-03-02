import json
import os
import shutil

from pathlib import Path

from odinpy.schemas import (
    CommandReadOutput,
    CommandReadOutputResponse,
    CommandTerminate,
    CommandTerminateResponse,
    CommandWriteInput,
    CommandWriteInputResponse,
    Error,
    ExecuteCommand,
    ExecuteCommandResponse,
    InstallNixPackage,
    InstallNixPackageResponse,
    NewTerminal,
    NewTerminalResponse,
    TerminalClose,
    TerminalCloseResponse,
    TerminalRead,
    TerminalReadResponse,
    TerminalWrite,
    TerminalWriteResponse,
    UninstallNixPackage,
    UninstallNixPackageResponse,
    UpsertFile,
)
from pydantic import BaseModel


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

    with schema_path.open("w") as f:
        json.dump(schema, f, indent=4)
        f.write("\n")  # Add newline at end of file


def generate():
    """Generate all schemas"""
    schemas_dir = Path(SCHEMAS_DIR)

    shutil.rmtree(schemas_dir)
    os.makedirs(schemas_dir, exist_ok=True)

    models = [
        Error,
        ExecuteCommand,
        ExecuteCommandResponse,
        CommandTerminate,
        CommandTerminateResponse,
        CommandReadOutput,
        CommandReadOutputResponse,
        CommandWriteInput,
        CommandWriteInputResponse,
        UpsertFile,
        NewTerminal,
        NewTerminalResponse,
        TerminalWrite,
        TerminalWriteResponse,
        TerminalRead,
        TerminalReadResponse,
        TerminalClose,
        TerminalCloseResponse,
        InstallNixPackage,
        InstallNixPackageResponse,
        UninstallNixPackage,
        UninstallNixPackageResponse,
    ]

    for model in models:
        write_schema(schemas_dir, model)
        print(f"Generated schema for {model.__name__}")


if __name__ == "__main__":
    generate()
