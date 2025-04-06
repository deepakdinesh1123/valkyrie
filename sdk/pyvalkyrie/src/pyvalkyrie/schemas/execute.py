from typing import List, Optional

from pydantic import BaseModel, Field


class EnvironmentVariable(BaseModel):
    name: str = Field(..., description="The name of the environment variable.")
    value: str = Field(..., description="The value of the environment variable.")


class ExecutionEnvironmentSpec(BaseModel):
    environment_variables: Optional[List[EnvironmentVariable]] = Field(
        None, description="A list of environment variables to set during execution."
    )
    languageDependencies: Optional[List[str]] = Field(
        None,
        description="A list of language-specific dependencies required for execution (e.g., Python packages, Node modules).",
    )
    systemDependencies: Optional[List[str]] = Field(
        None,
        description="A list of system-level dependencies (e.g., apt packages) needed for the execution environment.",
    )
    setup: Optional[str] = Field(
        None,
        description="Shell commands or script to run during setup, before code execution.",
    )


class ExecutionRequest(BaseModel):
    environment: Optional[ExecutionEnvironmentSpec] = Field(
        None,
        description="Specifications for the execution environment, including dependencies and setup instructions.",
    )
    code: str = Field(..., description="The code to execute.")
    language: str = Field(
        default="bash",
        description="The programming language of the code. Defaults to 'bash'.",
    )
    version: Optional[str] = Field(
        None, description="The version of the specified language to use."
    )
    max_retries: int = Field(
        default=5, description="Maximum number of retries in case of execution failure."
    )
    timeout: int = Field(
        default=-1, description="Execution timeout in seconds. -1 means no timeout."
    )
    cmdLineArgs: Optional[str] = Field(
        None, description="Command-line arguments to pass to the program."
    )
    compilerArgs: Optional[str] = Field(
        None, description="Arguments to pass to the compiler, if applicable."
    )
    command: Optional[str] = Field(
        None,
        description="The command to run to start execution, overriding default behavior.",
    )
    files: Optional[str] = Field(
        None, description="Additional files required for execution, encoded or inlined."
    )
    input: Optional[str] = Field(
        None,
        description="Standard input (stdin) to provide to the program during execution.",
    )


class ExecutionResponse(BaseModel):
    jobId: int
    events: str
    websocket: str


class ExecutionWSMessage(BaseModel):
    status: str
    jobId: int
    logs: Optional[str] = None
    errorMsg: Optional[str] = None
