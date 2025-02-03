from datetime import datetime

from typing import Optional, Literal

from pydantic import BaseModel, HttpUrl, WebsocketUrl

class Sandbox(BaseModel):
    sandboxId: int
    state: Literal['down', 'running', 'failed', 'stopped', 'creating', 'pending']
    URL: Optional[HttpUrl | WebsocketUrl] = None
    gitURL: Optional[HttpUrl] = None
    created_at: Optional[datetime] = None