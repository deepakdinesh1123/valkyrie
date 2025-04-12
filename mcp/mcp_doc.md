# Valkyrie MCP Server

The **Valkyrie MCP Server** provides a secure and flexible interface to execute code inside isolated containers. It supports multiple languages, dependency management, and runtime configurations.

## Tools

The MCP Server currently exposes the following tool:

### `execute_code`

This tool allows for secure code execution in a containerized environment, supporting multiple programming languages, compiler options, and system configurations.

---

## How to Use the Tools

### 1. `execute_code`

#### Description:
Executes the provided code in a sandboxed container environment. Supports specifying language, version, dependencies, files, input, and runtime configuration.

#### Arguments:

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `code` | `str` | ✅ | The actual code to execute. |
| `language` | `str` | ✅ | Language name (e.g. `python`, `go`, `rust`, `deno`). |
| `version` | `str` | ❌ | Full language version (e.g. `3.12.7`). Defaults to the latest stable if not specified. |
| `max_retries` | `int` | ❌ | Maximum number of retries on failure. |
| `timeout` | `int` | ❌ | Timeout in seconds. Default is `-1` (no timeout). |
| `cmd_line_args` | `str` | ❌ | Arguments passed to the program on execution. |
| `compiler_args` | `str` | ❌ | Arguments passed to the compiler (for compiled languages like `rust`, `go`, `c`, etc). |
| `command` | `str` | ❌ | Override the default command used to run the code. |
| `files` | `str` | ❌ | A JSON-encoded dictionary of filename-content pairs. Allows multi-file execution. |
| `input` | `str` | ❌ | Standard input to be passed to the program. |
| `environment_variables` | `List[Dict[str, str]]` | ❌ | List of environment variables. Each item must be a dictionary with `name` and `value`. |
| `language_dependencies` | `List[str]` | ❌ | Language-level dependencies (e.g. Python packages, ruby gems etc. Python packages are installed using UV and other packages can be installed only if nix ruspports them) |
| `system_dependencies` | `List[str]` | ❌ | System-level packages to be installed via the nix package manager (the name of the system dependency must be a valid nix package) |
| `setup` | `str` | ❌ | Additional shell commands to be executed before running the code (e.g. downloading data or compiling). |

---

## Example Usage

### Example 1: Run Python code

```json
{
  "tool": "execute_code",
  "args": {
    "language": "python",
    "code": "print('Hello from Valkyrie!')"
  }
}
```

### Example 1: Run Python code with requests dependency

```json
{
  "tool": "execute_code",
  "args": {
    "language": "import requests\nresp = requests.get('url').json()\nprint(resp)",
    "code": "print('Hello from Valkyrie!')",
    "language_dependencies": ["requests"]
  }
}
```
