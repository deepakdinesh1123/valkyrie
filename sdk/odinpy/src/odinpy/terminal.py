import json

from typing import Union
from uuid import UUID

from websocket import WebSocket, WebSocketTimeoutException

from .schemas import (
    Error,
    TerminalClose,
    TerminalCloseResponse,
    TerminalRead,
    TerminalReadResponse,
    TerminalWrite,
    TerminalWriteResponse,
)


class Terminal:
    def __init__(self, terminalId: UUID, agent: WebSocket):
        self._terminalId = terminalId
        self._agent = agent

    @property
    def terminalId(self) -> UUID:
        return self._terminalId

    def close_terminal(self) -> Union[TerminalCloseResponse, Error]:
        self._agent.send(TerminalClose(terminalId=self.terminalId).model_dump_json())
        try:
            resp = self._agent.recv()
        except WebSocketTimeoutException:
            raise TimeoutError("The server did not send a response in time.")

        message = json.loads(resp)

        try:
            terminalCloseResp = TerminalCloseResponse(**message)
            return terminalCloseResp
        except Exception:
            return Error(**message)

    def read_terminal(self) -> Union[TerminalReadResponse, Error]:
        self._agent.send(TerminalRead(terminalId=self.terminalId).model_dump_json())
        try:
            resp = self._agent.recv()
        except WebSocketTimeoutException:
            raise TimeoutError("The server did not send a response in time.")

        message = json.loads(resp)

        try:
            terminalReadResp = TerminalReadResponse(**message)
            return terminalReadResp
        except Exception:
            return Error(**message)

    def write_terminal(self, input: str) -> Union[TerminalWriteResponse, Error]:
        self._agent.send(
            TerminalWrite(terminalId=self.terminalId, input=input).model_dump_json()
        )
        try:
            resp = self._agent.recv()
        except WebSocketTimeoutException:
            raise TimeoutError("The server did not send a response in time.")

        message = json.loads(resp)

        try:
            terminalWriteResp = TerminalWriteResponse(**message)
            return terminalWriteResp
        except Exception:
            return Error(**message)
