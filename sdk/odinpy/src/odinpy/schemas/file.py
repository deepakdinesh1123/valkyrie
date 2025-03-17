from pydantic import BaseModel, Field
from typing import Optional

from .base import ResponseBase


class UpsertFile(BaseModel):
    """Represents a request to add or update a file in a sandbox."""

    msgType: str = Field(
        default="UpsertFile", frozen=True, description="Type of the message"
    )
    fileName: str = Field(description="Name of the file to be added or updated")
    path: str = Field(description="Path where the file should be created or updated")
    content: Optional[str] = Field(None, description="Content of the file")
    patch: Optional[str] = Field(None, description="Diff patch to apply to the file")


class DeleteFile(BaseModel):
    """Represents a request to delete a file in a sandbox."""

    msgType: str = Field(
        default="DeleteFile", frozen=True, description="Type of the message"
    )
    path: str = Field(description="Path of the file to be deleted")


class ReadFile(BaseModel):
    """Represents a request to read a file from a sandbox."""

    msgType: str = Field(
        default="ReadFile", frozen=True, description="Type of the message"
    )
    path: str = Field(description="Path of the file to be read")


class UpsertFileResponse(ResponseBase):
    """Represents a response to a file upsert request in a sandbox."""

    msgType: str = Field(
        default="UpsertFileResponse", frozen=True, description="Type of the message"
    )
    fileName: str = Field(description="Name of the file that was upserted")
    path: str = Field(description="Path where the file was upserted")


class DeleteFileResponse(ResponseBase):
    """Represents a response to a file delete request in a sandbox."""

    msgType: str = Field(
        default="DeleteFileResponse", frozen=True, description="Type of the message"
    )
    path: str = Field(description="Path of the file that was deleted")


class ReadFileResponse(ResponseBase):
    """Represents a response to a file read request in a sandbox."""

    msgType: str = Field(
        default="ReadFileResponse", frozen=True, description="Type of the message"
    )
    path: str = Field(description="Path of the file that was read")
    content: str = Field(description="Content of the file that was read")
