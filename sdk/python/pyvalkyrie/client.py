import requests

from .config import Config

class Client:

    def __init__(self, *args, **kwargs):
        self.config = Config()
    
    def create_sandbox(self) -> requests.Response:
        resp = requests.post(
            f"{self.config.ODIN_URL}/sandbox"
        )
        return resp 

    def check_sandbox_status(self):
        pass

    def new_sanbdox(self) -> int:
        pass