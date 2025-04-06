from typing import Optional, Tuple, Union

import requests

from pyvalkyrie.config import Config
from pyvalkyrie.exception import SandboxError
from pyvalkyrie.execution import Execution
from pyvalkyrie.log import setup_logging
from pyvalkyrie.sandbox import Sandbox
from pyvalkyrie.schemas import SandboxConfig
from pyvalkyrie.schemas.execute import ExecutionRequest
from requests.exceptions import RequestException


class Client:
    """
    Client for interacting with the sandboxes.

    This class provides methods to create and manage sandbox environments
    """

    def __init__(self, config: Config = Config()) -> None:
        """
        Initialize the sandbox client.

        Args:
            config (Optional[Config]): Configuration object. If None, creates a new Config instance.
        """
        self._config = config
        self.logger = setup_logging(config)
        self._session = requests.Session()
        self._session.headers.update(
            {
                "Content-Type": "application/json",
                "User-Agent": f"SandboxClient/{config.VERSION}",
            }
        )

    @property
    def config(self) -> Config:
        return self._config

    def _create_sandbox(
        self, sandbox_config: SandboxConfig
    ) -> Tuple[Union[int, str], bool]:
        """
        Send a request to create a new sandbox.

        Args:
            sandbox_config (SandboxConfig): Configuration for the new sandbox

        Returns:
            SandboxCreationResult: Object containing creation attempt results
        """
        try:
            url = f"{self.config.URL}/sandbox"
            response = self._session.post(url, json=sandbox_config.model_dump())
            response.raise_for_status()
            response_data = response.json()

            return response_data.get("sandboxId"), True
        except RequestException as e:
            return f"Request failed: {str(e)}", False
        except ValueError as e:
            return "Invalid JSON response from server", False

    def new_sandbox(
        self,
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
        sandboxConf = sandbox_config or SandboxConfig()
        result, success = self._create_sandbox(sandboxConf)

        if not success:
            raise SandboxError(f"Failed to create sandbox: {result}")

        return Sandbox(result, self.config, self.logger)

    def execute(self, execReq: ExecutionRequest):
        return Execution(execReq, self.logger, self.config)

    def __enter__(self) -> "Client":
        """Context manager entry."""
        return self

    def __exit__(self, exc_type, exc_val, exc_tb) -> None:
        """Context manager exit with proper cleanup."""
        self._session.close()
