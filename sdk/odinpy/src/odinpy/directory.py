import json
from typing import Union, Optional

from websocket import WebSocketTimeoutException

from .schemas import (
    UpsertDirectory,
    UpsertDirectoryResponse,
    DeleteDirectory,
    DeleteDirectoryResponse,
    ReadDirectory,
    ReadDirectoryResponse,
    Error,
)
from .log import logger


class Directory:
    def __init__(self, path: str, agent):
        """
        Initialize a Directory object for managing a directory in the sandbox.
        
        Args:
            path: Path of the directory in the sandbox
            agent: WebSocket connection to the sandbox agent
        """
        self._path = path
        self._agent = agent
    
    @property
    def path(self) -> str:
        return self._path
    
    def upsert(self, content: Optional[str] = None, patch: Optional[str] = None) -> Union[UpsertDirectoryResponse, Error]:
        """
        Update or create the directory in the sandbox.
        
        Args:
            content: Content of the directory (exclusive with patch)
            patch: Diff patch to apply to the directory (exclusive with content)
            
        Returns:
            UpsertDirectoryResponse on success, Error on failure
        """
        if content is None and patch is None:
            return Error(message="Either content or patch must be provided")
        
        payload = {
            "path": self._path
        }
        
        if content is not None:
            payload["content"] = content
        
        if patch is not None:
            payload["patch"] = patch
            
        request = UpsertDirectory(**payload)
        
        self._agent.send(request.model_dump_json())
        
        try:
            resp = self._agent.recv()
            message = json.loads(resp)
            
            try:
                return UpsertDirectoryResponse(**message)
            except Exception as e:
                logger.debug(f"Response from agent is {resp}")
                return Error(message=f"Failed to parse UpsertDirectoryResponse: {str(e)}")
                
        except WebSocketTimeoutException:
            return Error(message="WebSocket connection timed out while upserting directory.")
        except json.JSONDecodeError:
            logger.debug(f"Response from agent is {resp}")
            return Error(message="Failed to decode JSON response from the agent.")
        except Exception as e:
            return Error(message=f"An unexpected error occurred: {str(e)}")
    
    def delete(self) -> Union[DeleteDirectoryResponse, Error]:
        """
        Delete the directory from the sandbox.
        
        Returns:
            DeleteDirectoryResponse on success, Error on failure
        """
        request = DeleteDirectory(path=self._path)
        
        self._agent.send(request.model_dump_json())
        
        try:
            resp = self._agent.recv()
            message = json.loads(resp)
            
            try:
                return DeleteDirectoryResponse(**message)
            except Exception as e:
                logger.debug(f"Response from agent is {resp}")
                return Error(message=f"Failed to parse DeleteDirectoryResponse: {str(e)}")
                
        except WebSocketTimeoutException:
            return Error(message="WebSocket connection timed out while deleting directory.")
        except json.JSONDecodeError:
            logger.debug(f"Response from agent is {resp}")
            return Error(message="Failed to decode JSON response from the agent.")
        except Exception as e:
            return Error(message=f"An unexpected error occurred: {str(e)}")
    
    def read(self) -> Union[ReadDirectoryResponse, Error]:
        """
        Read the directory contents from the sandbox.
        
        Returns:
            ReadDirectoryResponse on success, Error on failure
        """
        request = ReadDirectory(path=self._path)
        
        self._agent.send(request.model_dump_json())
        
        try:
            resp = self._agent.recv()
            message = json.loads(resp)
            
            try:
                return ReadDirectoryResponse(**message)
            except Exception as e:
                logger.debug(f"Response from agent is {resp}")
                return Error(message=f"Failed to parse ReadDirectoryResponse: {str(e)}")
                
        except WebSocketTimeoutException:
            return Error(message="WebSocket connection timed out while reading directory.")
        except json.JSONDecodeError:
            logger.debug(f"Response from agent is {resp}")
            return Error(message="Failed to decode JSON response from the agent.")
        except Exception as e:
            return Error(message=f"An unexpected error occurred: {str(e)}") 