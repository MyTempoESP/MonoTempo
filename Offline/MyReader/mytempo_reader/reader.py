class ReaderOperationNotImplemented(NotImplementedError):
    def __init__(self, operation, reader, *args):
        self.message = f"{operation} Not implemented for this Reader ({type(reader)})"
        super(ReaderOperationNotImplemented, self).__init__(self.message, *args)


class BaseReader:
    def __init__(self, sender) -> None:
        self.sender = sender
        self.reader = None

    def finish_callback(self, *args) -> None: ...

    def inventory_start_callback(self, *args) -> None: ...

    # self, reader -- ignored --, tags
    def tag_report_callback(self):
        """Function to run each time the reader reports seeing tags."""
        ...

    def configure(self):
        raise ReaderOperationNotImplemented("Configure", self)

    def read(self):
        raise ReaderOperationNotImplemented("Read", self)

    def start(self):
        raise ReaderOperationNotImplemented("Start", self)
