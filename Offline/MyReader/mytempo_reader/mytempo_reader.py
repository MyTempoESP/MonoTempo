from os import getenv

from .chafon_reader import ChafonReader
from .impinj_reader import ImpinjReader
from .reader import BaseReader

ReaderType = BaseReader
if getenv("MYTEMPO_READER_TYPE") == "chafon":
    ReaderType = ChafonReader
else:
    ReaderType = ImpinjReader

# if called from cmdline
if __name__ == "__main__":
    reader: BaseReader = ReaderType()
    reader.start()


# exported
class Reader(ReaderType):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)
