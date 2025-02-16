import requests
import websocket

from typing import Union, Literal

from pydantic import HttpUrl, WebsocketUrl

from .terminal import Terminal
from .schemas import Sandbox as SandboxResponse, Error
from .config import Config

class Sandbox:
    
    def __init__(self, sandboxId: int, config: Config):
        self._sandboxId: int = sandboxId
        self._code_server: bool = True
        self._config: Config = config
        self.ws: websocket.WebSocket = None
        self.state: Literal["creating", "pending", "failed", "running"] = "pending"
    
    @property
    def config(self):
        return self._config
    
    @property
    def sandboxId(self):
        return self._sandboxId
    
    @property
    def sandboxURL(self):
        return self.sandboxURL
    
    @property
    def codeServer(self):
        return self._code_server
    
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

    def _connect_websocket(self) -> websocket.WebSocket:
        ws = websocket.WebSocket()
        ws.connect(self.sandboxURL, subprotocols=["sandbox"], timeout=20)
        return ws
    
    def _check_code_server(self):
        if self.codeServer:
            raise Exception(f"This sandbox has been initialised with code server, please visit the url {self.sandboxURL}")

    def new_terminal(self) -> Terminal:
        self._check_code_server()
        self.ws.send()