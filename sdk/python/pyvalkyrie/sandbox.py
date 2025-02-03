import websocket

from pydantic import HttpUrl, WebsocketUrl

from .terminal import Terminal
from .schemas import NewTerminal, NewTerminalResponse

class Sandbox:
    
    def __init__(self, sandboxId: int, sandboxURL: WebsocketUrl | HttpUrl):
        self._sandboxId = sandboxId
        self._sandboxURL = sandboxURL
        self._code_server = True
        if isinstance(sandboxURL, WebsocketUrl):
            self._code_server = False
            self.ws = self._connect_websocket()
    
    @property
    def sandboxId(self):
        return self._sandboxId
    
    @property
    def sandboxURL(self):
        return self.sandboxURL
    
    @property
    def codeServer(self):
        return self._code_server

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