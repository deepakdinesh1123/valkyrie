from typing import Dict, List

from mcp.server.fastmcp import FastMCP
from pyvalkyrie import Client
from pyvalkyrie.config import Config
from pyvalkyrie.schemas.execute import ExecutionRequest


mcp = FastMCP("Valkyrie")
config = Config()
cli = Client(config)


@mcp.tool(
    name="Valkyrie",
    description="A tool for executing code in containerized environments using nix",
)
async def execute_code(
    code: str,
    language: str = "bash",
    version: str = None,
    max_retries: int = 5,
    timeout: int = -1,
    cmd_line_args: str = None,
    compiler_args: str = None,
    command: str = None,
    files: str = None,
    input: str = None,
    environment_variables: List[Dict[str, str]] = None,
    language_dependencies: List[str] = None,
    system_dependencies: List[str] = None,
    setup: str = None,
):
    # Construct environment spec if any environment-related parameters are provided
    environment = None
    if any([environment_variables, language_dependencies, system_dependencies, setup]):
        environment = {
            "environment_variables": environment_variables,
            "languageDependencies": language_dependencies,
            "systemDependencies": system_dependencies,
            "setup": setup,
        }

    # Construct execution request
    exec_req = ExecutionRequest(
        code=code,
        language=language,
        version=version,
        max_retries=max_retries,
        timeout=timeout,
        cmdLineArgs=cmd_line_args,
        compilerArgs=compiler_args,
        command=command,
        files=files,
        input=input,
        environment=environment,
    )

    res = cli.execute(exec_req)
    if res.status == "completed":
        return res.output
    else:
        return f"Execution failed: {res.output}"


if __name__ == "__main__":
    mcp.run(transport="stdio")
