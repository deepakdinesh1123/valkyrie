from pydantic import BaseModel, Field


class UpsertFile(BaseModel):
    """Represents a request to add or update a file in a sandbox."""

    msgType: str = Field(
        default="UpsertFile", frozen=True, description="Type of the message"
    )
    sandboxId: int = Field(
        description="ID of the sandbox in which to add or update the file"
    )
    fileName: str = Field(description="Name of the file to be added or updated")
    path: str = Field(description="Path where the file should be created or updated")
    content: str = Field(..., description="Content of the file")
    patch: str = Field(..., description="Diff patch to apply to the file")


class DeleteFile(BaseModel):
    """Represents a request to delete a file in a sandbox."""

    msgType: str = Field(
        default="DeleteFile", frozen=True, description="Type of the message"
    )
    sandboxId: int = Field(
        description="ID of the sandbox from which to delete the file"
    )
    path: str = Field(description="Path of the file to be deleted")


class ReadFile(BaseModel):
    """Represents a request to read a file from a sandbox."""

    msgType: str = Field(
        default="ReadFile", frozen=True, description="Type of the message"
    )
    sandboxId: int = Field(description="ID of the sandbox from which to read the file")
    path: str = Field(description="Path of the file to be read")
