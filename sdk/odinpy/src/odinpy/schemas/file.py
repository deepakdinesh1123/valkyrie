from pydantic import BaseModel, Field


class AddFile(BaseModel):
    """Add a file to the sandbox"""

    msgType: str = Field(default="AddFile", frozen=True)
    sandboxId: int = Field(..., description="ID of the sandbox to add the file to")
    fileName: str = Field(..., description="Name of the file")
    path: str = Field(..., description="Path where to create the file")
    content: str = Field(..., description="File content")
