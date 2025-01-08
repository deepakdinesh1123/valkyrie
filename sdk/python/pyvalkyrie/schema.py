# /// script
# requires-python = ">=3.12"
# dependencies = [
#     "pydantic",
# ]
# ///
import os
import json
from typing import List

from pydantic import BaseModel

SCHEMAS_DIR = "schemas/agent"

def write_schema(file: str, schema: str):
    with open(os.path.join(SCHEMAS_DIR, file), "w") as f:
        f.write(json.dumps(schema, indent=4))

class EnvironmentVariable(BaseModel):
    key: str
    value: str

class ExecuteCommand(BaseModel):
    sandboxId: int
    command: str
    stdin: bool = False
    stderr: bool = True
    stdout: bool = True
    workDir: str = ""
    env: List[str] = []

if __name__ == "__main__":
    write_schema("execute_command.json", ExecuteCommand.model_json_schema())
