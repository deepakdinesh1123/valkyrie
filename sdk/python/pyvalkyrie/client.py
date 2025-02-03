import requests
from time import time, sleep
from typing import Union, Tuple, Optional

from .config import Config
from .schemas import Sandbox as SandboxResponse, Error
from .sandbox import Sandbox

class Client:
    """
    Client for interacting with the sandboxes.
    """
    def __init__(self) -> None:
        """
        Initializes the client with configuration settings.
        """
        self.config = Config()
    
    def _create_sandbox(self) -> Tuple[str, Optional[int]]:
        """
        Sends a request to create a new sandbox.
        
        Returns:
            Tuple[Optional[str], Optional[int]]: A tuple containing an error message (if any) 
            and the sandbox ID (if creation was successful).
        """
        resp = requests.post(url=f"{self.config.ODIN_URL}/sandbox")
        response = resp.json()
        
        if resp.status_code != 200:
            return response.get("message"), None
        
        return response.get("result"), response.get("sandboxId")

    def check_sandbox_status(self, sandbox_id: int) -> Union[SandboxResponse, Error]:
        """
        Checks the status of a sandbox given its ID.
        
        Args:
            sandbox_id (int): The ID of the sandbox to check.
        
        Returns:
            Union[SandboxResponse, Error]: The sandbox status if successful, otherwise an error.
        """
        resp = requests.get(url=f"{self.config.ODIN_URL}/{sandbox_id}")
        
        if resp.status_code != 200:
            return Error(**resp.json())
        
        return SandboxResponse(**resp.json())

    def new_sandbox(self) -> Sandbox:
        """
        Creates a new sandbox and waits until it is running.
        
        Returns:
            Sandbox: An instance of the Sandbox class with the sandbox details.
        
        Raises:
            Exception: If the sandbox could not be created or an error occurs during status checks.
        """
        start_time = time()
        
        error_message, sandbox_id = self._create_sandbox()
        
        if not sandbox_id:
            raise Exception(f"Could not create sandbox: {error_message}")
        
        while time() - start_time < self.config.SANDBOX_CREATION_TIMEOUT:
            status_response = self.check_sandbox_status(sandbox_id)
            
            if isinstance(status_response, Error):
                raise Exception(f"Error checking sandbox status: {status_response.message}")
            
            if isinstance(status_response, SandboxResponse):
                if status_response.state == "running":
                    return Sandbox(sandboxId=status_response.sandboxId, sandboxURL=status_response.URL)
                elif status_response.state == "pending":
                    pass  # Continue waiting
            
            sleep(2)
        
        raise Exception("Sandbox creation timed out")
