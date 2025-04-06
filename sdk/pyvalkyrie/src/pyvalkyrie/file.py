import json
import logging

from typing import Optional, Union

from pyvalkyrie.schemas import (
    DeleteFile,
    DeleteFileResponse,
    Error,
    ReadFile,
    ReadFileResponse,
    UpsertFile,
    UpsertFileResponse,
)
from websocket import WebSocketTimeoutException


class File:
    def __init__(self, path: str, agent, logger: logging.Logger):
        """
        Initialize a File object for managing a file in the sandbox.

        Args:
            path: Path of the file in the sandbox
            agent: WebSocket connection to the sandbox agent
        """
        self._path = path
        self._agent = agent
        self.logger = logger

    @property
    def path(self) -> str:
        return self._path

    def upsert(
        self, content: Optional[str] = None, patch: Optional[str] = None
    ) -> Union[UpsertFileResponse, Error]:
        """
        Update or create the file in the sandbox.

        Args:
            content: Content of the file (exclusive with patch)
            patch: Diff patch to apply to the file (exclusive with content)

        Returns:
            UpsertFileResponse on success, Error on failure
        """
        if content is None and patch is None:
            return Error(message="Either content or patch must be provided")

        filename = self._path.split("/")[-1]
        payload = {"fileName": filename, "path": self._path}

        if content is not None:
            payload["content"] = content

        if patch is not None:
            payload["patch"] = patch

        request = UpsertFile(**payload)

        self._agent.send(request.model_dump_json())

        try:
            resp = self._agent.recv()
            message = json.loads(resp)

            try:
                return UpsertFileResponse(**message)
            except Exception as e:
                self.logger.debug(f"Response from agent is {resp}")
                return Error(message=f"Failed to parse UpsertFileResponse: {str(e)}")

        except WebSocketTimeoutException:
            return Error(message="WebSocket connection timed out while upserting file.")
        except json.JSONDecodeError:
            self.logger.debug(f"Response from agent is {resp}")
            return Error(message="Failed to decode JSON response from the agent.")
        except Exception as e:
            return Error(message=f"An unexpected error occurred: {str(e)}")

    def delete(self) -> Union[DeleteFileResponse, Error]:
        """
        Delete the file from the sandbox.

        Returns:
            DeleteFileResponse on success, Error on failure
        """
        request = DeleteFile(path=self._path)

        self._agent.send(request.model_dump_json())

        try:
            resp = self._agent.recv()
            message = json.loads(resp)

            try:
                return DeleteFileResponse(**message)
            except Exception as e:
                self.logger.debug(f"Response from agent is {resp}")
                return Error(message=f"Failed to parse DeleteFileResponse: {str(e)}")

        except WebSocketTimeoutException:
            return Error(message="WebSocket connection timed out while deleting file.")
        except json.JSONDecodeError:
            self.logger.debug(f"Response from agent is {resp}")
            return Error(message="Failed to decode JSON response from the agent.")
        except Exception as e:
            return Error(message=f"An unexpected error occurred: {str(e)}")

    def read(self) -> Union[ReadFileResponse, Error]:
        """
        Read the file content from the sandbox.

        Returns:
            ReadFileResponse on success, Error on failure
        """
        request = ReadFile(path=self._path)

        self._agent.send(request.model_dump_json())

        try:
            resp = self._agent.recv()
            message = json.loads(resp)

            try:
                return ReadFileResponse(**message)
            except Exception as e:
                self.logger.debug(f"Response from agent is {resp}")
                return Error(message=f"Failed to parse ReadFileResponse: {str(e)}")

        except WebSocketTimeoutException:
            return Error(message="WebSocket connection timed out while reading file.")
        except json.JSONDecodeError:
            self.logger.debug(f"Response from agent is {resp}")
            return Error(message="Failed to decode JSON response from the agent.")
        except Exception as e:
            return Error(message=f"An unexpected error occurred: {str(e)}")
