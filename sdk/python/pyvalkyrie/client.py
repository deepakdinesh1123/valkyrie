import requests
from typing import Optional
from requests.exceptions import RequestException
from dataclasses import dataclass

from .config import Config
from .schemas import SandboxConfig
from .sandbox import Sandbox

@dataclass
class SandboxCreationResult:
    sandbox_id: Optional[int] = None
    error_message: Optional[str] = None
    success: bool = False

class SandboxError(Exception):
    """Custom exception for sandbox-related errors."""
    pass

class Client:
    """
    Client for interacting with the sandboxes.
    
    This class provides methods to create and manage sandbox environments
    """
    
    def __init__(self) -> None:
        """
        Initialize the sandbox client.
        
        Args:
            config (Optional[Config]): Configuration object. If None, creates a new Config instance.
        """
        self._config = Config()
        self._session = requests.Session()
        self._session.headers.update({
            "Content-Type": "application/json",
            "User-Agent": f"SandboxClient/{self._config.VERSION}"
        })
    
    @property
    def config(self) -> Config:
        """Get the current configuration."""
        return self._config
    
    def _create_sandbox(self, sandbox_config: SandboxConfig) -> SandboxCreationResult:
        """
        Send a request to create a new sandbox.
        
        Args:
            sandbox_config (SandboxConfig): Configuration for the new sandbox
            
        Returns:
            SandboxCreationResult: Object containing creation attempt results
        """
        try:
            url = f"{self.config.ODIN_URL}/sandbox"
            response = self._session.post(
                url,
                json=sandbox_config.model_dump()
            )
            response.raise_for_status()
            response_data = response.json()
            
            return SandboxCreationResult(
                sandbox_id=response_data.get("sandboxId"),
                error_message=response_data.get("message"),
                success=bool(response_data.get("sandboxId"))
            )
        except RequestException as e:
            return SandboxCreationResult(
                error_message=f"Request failed: {str(e)}",
                success=False
            )
        except ValueError as e:
            return SandboxCreationResult(
                error_message="Invalid JSON response from server",
                success=False
            )

    def new_sandbox(self, 
        sandbox_config: Optional[SandboxConfig] = None,
    ) -> Sandbox:
        """
        Create a new sandbox and wait until it is running.
        
        Args:
            sandbox_config (Optional[SandboxConfig]): Configuration for the new sandbox.
                If None, uses default configuration.
        
        Returns:
            Sandbox: An instance of the Sandbox class with the sandbox details
        
        Raises:
            SandboxError: If the sandbox could not be created or an error occurs
        """
        config = sandbox_config or SandboxConfig()
        result = self._create_sandbox(config)
        
        if not result.success:
            raise SandboxError(f"Failed to create sandbox: {result.error_message}")
            
        return Sandbox(self, result.sandbox_id)
    
    def __enter__(self) -> 'Client':
        """Context manager entry."""
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb) -> None:
        """Context manager exit with proper cleanup."""
        self._session.close()