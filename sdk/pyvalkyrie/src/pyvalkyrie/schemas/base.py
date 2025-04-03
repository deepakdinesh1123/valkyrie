from pydantic import BaseModel


class ResponseBase(BaseModel):
    success: bool
    msg: str
