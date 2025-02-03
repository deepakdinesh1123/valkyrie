from uuid import UUID

class Terminal:

    def __init__(self, terminalId: UUID):
        self._terminalId = terminalId
    
    @property
    def terminal(self):
        return self._terminalId
    
    def close_terminal(self):
        pass

    def read_terminal(self):
        pass

    def write_terminal(self, input: str):
        pass