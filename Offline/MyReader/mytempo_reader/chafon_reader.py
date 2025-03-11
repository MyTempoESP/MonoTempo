import logging
import sys
from typing import Never

import pendulum

from .epc import EPC, InvalidEPCError
from .reader import BaseReader
from .tag import Tag

logging.basicConfig(level=logging.INFO)

chafon_imported, err = False, None

try:
    import pychafon

    chafon_imported = True
except ImportError as e:
    logging.info("Chafon not imported!")
    err = e


class _ChafonReader(BaseReader):
    def tag_report_callback(self, code):
        # print(code)

        try:
            epc = EPC(code)
        except InvalidEPCError:
            logging.info(f"Ignoring invalid tag: '{code}'...")
            return

        try:
            last_seen = pendulum.now().in_timezone("America/Sao_Paulo")  # Damn.
            tag_obj: Tag = Tag(epc, last_seen)
            self.sender.send_tag(tag_obj)
        except Exception as e:
            logging.warn(f"Ignoring error at tag callback '{e}'...")

    def configure(self) -> bool:
        self.reader = pychafon

        try:
            self.reader.open()
        except RuntimeError:
            return False

        self.reader.setTagCallback(self.tag_report_callback)

        return True

    def read(self) -> Never:
        try:
            self.reader.read(0)
        except (SystemExit, KeyboardInterrupt):
            logging.info("Stopping readers...")

            try:
                self.reader.close()
            except Exception:
                logging.warning("Ignoring exception at Disconnect")

    def start(self) -> Never:
        if not self.reader:
            if not self.configure():
                logging.warning("No Reader connected!")

                # XXX:
                #   Exit if reader not connected ( USE WITH RESTART: ALWAYS on docker compose )
                #   this is required in chafon readers because for some wacky reason my C++ program
                #   or their .so file doesn't connect sometimes if you loop it, unless the lib is
                #   fully re-loaded.
                sys.exit(1)

        logging.info("Reader connected!")
        self.read()


class ChafonReader(_ChafonReader):
    def __init__(self, *args, **kwargs):
        if chafon_imported:
            super().__init__(*args, **kwargs)
        else:
            raise ImportError(f"Couldn't import pychafon! Exception: {err}")
