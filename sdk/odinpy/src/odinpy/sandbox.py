import json
import time

from time import sleep
from typing import Literal, Optional, Union

import requests


from websocket import WebSocket, WebSocketException, WebSocketTimeoutException

from .config import config
from .exception import SandboxError, SandboxTimeoutError
from .log import logger
from .schemas import Error, NewTerminal, NewTerminalResponse
from .schemas import Sandbox as SandboxResponse
from .terminal import Terminal
from .file import File
from .directory import Directory
from .websocket import websocket_connection

from .schemas import (
    UpsertFileResponse,
    DeleteFileResponse,
    ReadFileResponse,
    UpsertDirectoryResponse,
    DeleteDirectoryResponse,    
    ReadDirectoryResponse,
)
class Sandbox:
    def __init__(self, sandboxId: int):
        self._sandboxId: int = sandboxId
        self._agent: Optional[WebSocket] = None
        self._state: Literal["creating", "pending", "failed", "running"] = "pending"
        self._sandboxURL: str = None
        self._sandboxAgentURL: str = None

        res = self._check_sandbox_creation_result()
        if isinstance(res, SandboxResponse):
            self._set_sandboxAgentURL(res.agentURL)
            self._set_state(res.state)
            self._set_sandboxURL(res.URL)
            self._set_agent(self._connect_agent())
        else:
            raise SandboxError(f"Sandbox creation failed: {res.message}")

    @property
    def sandboxId(self) -> int:
        return self._sandboxId

    @property
    def sandboxURL(self):
        return self._sandboxURL

    def _set_sandboxURL(self, url: str) -> None:
        self._sandboxURL = url

    @property
    def sandboxAgentURL(self) -> Optional[str]:
        return self._sandboxAgentURL

    def _set_sandboxAgentURL(self, url: str) -> None:
        self._sandboxAgentURL = url

    @property
    def state(self) -> Literal["creating", "pending", "failed", "running"]:
        return self._state

    def _set_state(
        self, state: Literal["creating", "pending", "failed", "running"]
    ) -> None:
        self._state = state

    @property
    def agent(self) -> Optional[WebSocket]:
        return self._agent

    def _set_agent(self, agent: Optional[WebSocket]) -> None:
        self._agent = agent

    def check_sandbox_status(self) -> Union[SandboxResponse, Error]:
        """
        Checks the status of a sandbox given its ID.

        Args:
            sandbox_id (int): The ID of the sandbox to check.

        Returns:
            Union[SandboxResponse, Error]: The sandbox status if successful, otherwise an error.
        """
        url = f"{config.ODIN_URL}/{self.sandboxId}"

        try:
            resp = requests.get(url)
            resp.raise_for_status()
            return SandboxResponse(**resp.json())

        except requests.RequestException as e:
            return Error(message=f"Request failed: {str(e)}")

        except ValueError:
            return Error(message="Invalid JSON response from server")

    def _check_sandbox_creation_result(self) -> Union[SandboxResponse, Error]:
        """
        Check the result of sandbox creation via WebSocket connection.

        Args:
            timeout: Maximum time in seconds to wait for sandbox creation

        Returns:
            SandboxResponse on success, Error on failure

        Raises:
            SandboxTimeoutError: If sandbox creation exceeds timeout
            Error: For other sandbox errors
        """
        sandbox_status_url = f"{'wss' if config.IS_ODIN_SECURE else 'ws'}://{config.ODIN_HOST}/sandboxes/{self.sandboxId}/status/ws"
        start_time = time.time()

        try:
            with websocket_connection(sandbox_status_url) as ws:
                logger.debug("Connected to websocket")

                while time.time() - start_time < config.SANDBOX_CREATION_TIMEOUT:
                    try:
                        resp = ws.recv()
                        if not resp:
                            logger.debug("No response received")
                            sleep(2)
                            continue

                        message = json.loads(resp)
                        logger.debug(f"Received websocket message: {message}")

                        event = message.get("event")
                        if not event:
                            logger.warning("No event specified in message")
                            continue

                        data = message.get("data", {})
                        logger.debug(f"Event: {event}, Data: {data}")

                        if event == "status":
                            state = data.get("state")
                            if not state:
                                logger.warning("No state in status event")
                                continue

                            self._state = state

                            if state in ["creating", "pending"]:
                                logger.debug(f"Sandbox Status: {state}")
                            elif state == "failed":
                                logger.error(f"Sandbox creation failed: {data}")
                                return Error(message=f"Sandbox creation failed: {data}")
                            elif state == "running":
                                return SandboxResponse(**data)

                        elif event == "error":
                            error_msg = f"Received error event: {data}"
                            logger.error(error_msg)
                            return Error(message=error_msg)

                        elif event == "close":
                            error_msg = f"Connection closed unexpectedly: {data}"
                            logger.error(error_msg)
                            return Error(message=error_msg)

                    except json.JSONDecodeError as e:
                        logger.error(f"Failed to decode response: {e}")
                        continue

                    except WebSocketException as e:
                        logger.error(f"WebSocket error: {e}")
                        return Error(message=f"WebSocket error: {e}")

                    except Exception as e:
                        logger.error(f"Unexpected error: {e}", exc_info=True)
                        return Error(message=f"Unexpected error: {e}")

                raise SandboxTimeoutError(
                    f"Sandbox creation timed out after {config.SANDBOX_CREATION_TIMEOUT} seconds"
                )

        except WebSocketException as e:
            error_msg = f"Failed to connect to WebSocket: {e}"
            logger.error(error_msg)
            return Error(message=error_msg)

    def _connect_agent(self) -> WebSocket:
        ws = WebSocket()
        logger.debug(f"Websocket URL is {self.sandboxAgentURL}")
        ws.connect(self.sandboxAgentURL, subprotocols=["sandbox"], timeout=20)
        ws.settimeout(config.SANDBOX_AGENT_TIMEOUT)
        return ws

    def new_terminal(self) -> Union[Terminal, Error]:
        if self.state != "running":
            return Error(
                message=f"Sandbox is not running, current state is: {self.state}"
            )

        self.agent.send(payload=NewTerminal().model_dump_json())

        try:
            resp = self.agent.recv()
            message = NewTerminalResponse.model_validate_json(resp)

            term = Terminal(terminalId=message.terminalID, agent=self.agent)
            return term
        except WebSocketTimeoutException:
            return Error(
                message="WebSocket connection timed out while creating a new terminal."
            )
        except json.JSONDecodeError:
            logger.debug(f"Response from agent is {resp}")
            return Error(message="Failed to decode JSON response from the agent.")
        except Exception as e:
            return Error(message=f"An unexpected error occurred: {str(e)}")
        

    def get_file(self, path: str) -> File:
        """
        Get a File object for the specified path.

        Args:
            path: Path of the file in the sandbox

        Returns:
            File: A File object for the specified path
        """
        if self.state != "running":
            return Error(
                message=f"Sandbox is not running, current state is: {self.state}"
            )

        return File(path=path, agent=self.agent)
    
    def upsert_file(self, path: str, content: Optional[str] = None, patch: Optional[str] = None) -> Union[UpsertFileResponse, Error]:
        """
        Update or create a file in the sandbox.
     
        Args:
            path: Path where the file should be created or updated
            content: Content of the file (exclusive with patch)
            patch: Diff patch to apply to the file (exclusive with content)
            
        Returns:
            UpsertFileResponse on success, Error on failure
        """
        if content is None and patch is None:
            return Error(message="Either content or patch must be provided")
     
        file = self.get_file(path)
        if isinstance(file, Error):
            return file
     
        return file.upsert(content=content, patch=patch)
    
    def delete_file(self, path: str) -> Union[DeleteFileResponse, Error]:
        """
        Delete a file from the sandbox.

        Args:
            path: Path of the file to be deleted

        Returns:
            DeleteFileResponse on success, Error on failure
        """
        file = self.get_file(path)
        if isinstance(file, Error):
            return file

        return file.delete()

    def read_file(self, path: str) -> Union[ReadFileResponse, Error]:
        """
        Read a file from the sandbox.

        Args:
            path: Path of the file to be read

        Returns:
            ReadFileResponse on success, Error on failure
        """
        file = self.get_file(path)
        if isinstance(file, Error):
            return file

        return file.read()
    

    def get_directory(self, path: str) -> Union[Directory, Error]:
        """
        Get a Directory object for managing a directory in the sandbox.
        
        Args:
            path: Path of the directory in the sandbox
            
        Returns:
            Directory object on success, Error on failure
        """        
        if self.state != "running":
            return Error(
                message=f"Sandbox is not running, current state is: {self.state}"
            )

        return Directory(path=path, agent=self.agent)
    
    def upsert_directory(self, path: str, content: Optional[str] = None, patch: Optional[str] = None) -> Union[UpsertDirectoryResponse, Error]:
        """
        Update or create a directory in the sandbox.
     
        Args:
            path: Path where the directory should be created or updated
            content: Content of the directory (exclusive with patch)
            patch: Diff patch to apply to the directory (exclusive with content)
            
        Returns:
            UpsertDirectoryResponse on success, Error on failure
        """
        if content is None and patch is None:
            return Error(message="Either content or patch must be provided")
     
        directory = self.get_directory(path)
        if isinstance(directory, Error):
            return directory
     
        return directory.upsert(content=content, patch=patch)
    
    def delete_directory(self, path: str) -> Union[DeleteDirectoryResponse, Error]:
        """
        Delete a directory from the sandbox.

        Args:
            path: Path of the directory to be deleted

        Returns:
            DeleteDirectoryResponse on success, Error on failure
        """
        directory = self.get_directory(path)
        if isinstance(directory, Error):
            return directory

        return directory.delete()

    def read_directory(self, path: str) -> Union[ReadDirectoryResponse, Error]:
        """
        Read a directory from the sandbox.

        Args:
            path: Path of the directory to be read

        Returns:
            ReadDirectoryResponse on success, Error on failure
        """
        directory = self.get_directory(path)
        if isinstance(directory, Error):
            return directory

        return directory.read() 
    