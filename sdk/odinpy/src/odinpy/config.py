from typing import Literal
from urllib.parse import urlparse

from pydantic import Field, HttpUrl
from pydantic_settings import BaseSettings


class Config(BaseSettings):
    ODIN_URL: HttpUrl = Field(default="http://localhost:8080/api")
    SANDBOX_CREATION_TIMEOUT: int = Field(default=300)
    SANDBOX_AGENT_TIMEOUT: int = Field(default=300)
    VERSION: str = Field(default="0.0.1")
    LOG_LEVEL: Literal["DEBUG", "INFO", "ERROR", "WARNING"] = Field(default="DEBUG")
    # ODIN_USER_TOKEN: str
    # ODIN_ADMIN_TOKEN: str

    @property
    def ODIN_HOST(self) -> str:
        return urlparse(str(self.ODIN_URL)).netloc

    @property
    def IS_ODIN_SECURE(self) -> bool:
        return urlparse(str(self.ODIN_URL)).scheme == "https"


config = Config()
