from pydantic import BaseModel, Field

from .base import ResponseBase


class UpsertDirectory(BaseModel):
    """Represents a request to add or update a directory in a sandbox."""

    msgType: str = Field(
        default="UpsertDirectory", frozen=True, description="Type of the message"
    )
    path: str = Field(
        description="Path where the directory should be created or updated"
    )


class DeleteDirectory(BaseModel):
    """Represents a request to delete a directory in a sandbox."""

    msgType: str = Field(
        default="DeleteDirectory", frozen=True, description="Type of the message"
    )
    path: str = Field(description="Path of the directory to be deleted")


class ReadDirectory(BaseModel):
    """Represents a request to read a directory from a sandbox."""

    msgType: str = Field(
        default="ReadDirectory", frozen=True, description="Type of the message"
    )
    path: str = Field(description="Path of the directory to be read")


class UpsertDirectoryResponse(ResponseBase):
    """Represents a response to a directory upsert request in a sandbox."""

    msgType: str = Field(
        default="UpsertDirectoryResponse",
        frozen=True,
        description="Type of the message",
    )
    path: str = Field(description="Path where the directory was upserted")


class DeleteDirectoryResponse(ResponseBase):
    """Represents a response to a directory delete request in a sandbox."""

    msgType: str = Field(
        default="DeleteDirectoryResponse",
        frozen=True,
        description="Type of the message",
    )
    path: str = Field(description="Path of the directory that was deleted")


class ReadDirectoryResponse(ResponseBase):
    """Represents a response to a directory read request in a sandbox."""

    msgType: str = Field(
        default="ReadDirectoryResponse", frozen=True, description="Type of the message"
    )
    path: str = Field(description="Path of the directory that was read")
    contents: str = Field(description="Content of the directory that was read")
