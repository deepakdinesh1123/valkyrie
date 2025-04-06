import json
import logging

from time import sleep
from typing import Literal

import requests

from pyvalkyrie.config import Config
from pyvalkyrie.schemas.execute import (
    ExecutionRequest,
    ExecutionResponse,
    ExecutionWSMessage,
)
from pyvalkyrie.websocket import websocket_connection


class Execution:
    def __init__(
        self, execReq: ExecutionRequest, logger: logging.Logger, config: Config
    ):
        self._execReq = execReq
        self._status: Literal[
            "pending", "scheduled", "completed", "failed", "cancelled"
        ] = "pending"
        self.logger = logger
        self._config = config
        self._output = ""
        self._execId = ""
        self.execute()

    @property
    def execReq(self) -> ExecutionRequest:
        return self._execReq

    @property
    def status(
        self,
    ) -> Literal["pending", "scheduled", "completed", "failed", "cancelled"]:
        return self._status

    @property
    def output(self):
        return self._output

    @property
    def config(self) -> Config:
        return self._config

    @property
    def execId(self) -> str:
        return self._execId

    def execute(self):
        # Set up the execution request
        url = f"{self.config.URL}/executions/execute"
        headers = {"content-type": "application/json"}

        try:
            # Send execution request
            resp = requests.post(
                url=url,
                data=self.execReq.model_dump_json(exclude_none=True),
                headers=headers,
            )
            try:
                resp.raise_for_status()
            except requests.Request:
                return "failure", resp.json()

            # Parse response
            execRes = ExecutionResponse(**resp.json())
            self._execId = execRes.jobId  # Save the execution ID from the response

            # Set status to scheduled once we have a valid execution ID
            self._status = "scheduled"

            # Build websocket URL
            websocket_url = f"{'wss' if self.config.IS_SECURE else 'ws'}://{self.config.HOST}{execRes.websocket}"
            self.logger.info(f"Connecting to websocket at {websocket_url}")

            # Connect to websocket and monitor execution
            with websocket_connection(websocket_url) as ws:
                while True:
                    self.logger.info("Waiting for execution updates...")
                    message = ws.recv()
                    message_data = ExecutionWSMessage(**json.loads(message))

                    self.logger.info(
                        f"Current status: {self._status}, Received websocket message: {message_data}"
                    )

                    # Update status and output
                    self._status = message_data.status

                    # Handle different message types
                    if message_data.logs:
                        self._output = message_data.logs
                    elif message_data.errorMsg:
                        self._output = message_data.errorMsg

                    # Exit the loop if we've reached a terminal state
                    if self._status in ["completed", "failed", "cancelled"]:
                        break

                    sleep(4)

        except requests.exceptions.RequestException as e:
            self.logger.error(f"HTTP request failed: {str(e)}")
            self._status = "failed"
            self._output = f"Request failed: {str(e)}"
        except json.JSONDecodeError as e:
            self.logger.error(f"Failed to parse JSON response: {str(e)}")
            self._status = "failed"
            self._output = f"JSON parse error: {str(e)}"
        except Exception as e:
            self.logger.error(f"Error during execution: {str(e)}")
            self._status = "failed"
            self._output = f"Execution error: {str(e)}"

        self.logger.info(f"Execution completed with status: {self._status}")
        return self._status, self._output
