import json
import logging

from uuid import UUID

from pyvalkyrie.schemas import (
    TerminalClose,
    TerminalCloseResponse,
    TerminalRead,
    TerminalReadResponse,
    TerminalWrite,
    TerminalWriteResponse,
)
from websocket import WebSocket, WebSocketTimeoutException


class Terminal:
    def __init__(self, terminalId: UUID, agent: WebSocket, logger: logging.Logger):
        self._terminalId = terminalId
        self._agent = agent
        self.logger = logger

    @property
    def terminalId(self) -> UUID:
        return self._terminalId

    def close_terminal(self) -> TerminalCloseResponse:
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
            return TerminalCloseResponse(**message)

    def read_terminal(self) -> TerminalReadResponse:
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
            return TerminalReadResponse(**message)

    def write_terminal(self, input: str) -> TerminalWriteResponse:
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
            return TerminalWriteResponse(**message)
