from pydantic_settings import BaseSettings
from pydantic import AnyHttpUrl, Field

class Config(BaseSettings):
    ODIN_URL: AnyHttpUrl = Field(default="http://localhost:8080/api")
    SANDBOX_CREATION_TIMEOUT: int = Field(default=60) 
    # ODIN_USER_TOKEN: str
    # ODIN_ADMIN_TOKEN: str