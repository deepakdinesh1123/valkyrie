from datetime import datetime
from typing import List, Literal, Optional

from pydantic import BaseModel


class Sandbox(BaseModel):
    sandboxId: int
    state: Literal["down", "running", "failed", "stopped", "creating", "pending"]
    URL: str = None
    agentURL: str = None
    created_at: Optional[datetime] = None


class SandboxConfig(BaseModel):
    nix_flake: str = ""
    languages: List[str] = []
    system_dependencies: List[str] = []
    shell_hook: str = ""
    # devenv: str = Field(default_factory=get_default_devenv_config)
