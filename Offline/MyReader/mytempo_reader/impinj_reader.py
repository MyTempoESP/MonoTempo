import logging
from time import sleep
from typing import Never

import pendulum
from sllurp.llrp import LLRPReaderClient, LLRPReaderConfig, LLRPReaderState

from .epc import EPC, InvalidEPCError
from .reader import BaseReader
from .tag import Tag

logging.basicConfig(level=logging.INFO)


class ImpinjReader(BaseReader):
    def tag_report_callback(self, _, tags: list):
        if len(tags) == 0:
            return

        for tag in tags:
            try:
                epc = EPC.from_tag_data(tag)
            except InvalidEPCError:
                logging.info(
                    "Ignoring invalid tag...",
                )
                continue

            # first_seen = tag["FirstSeenTimestampUTC"]
            # last_seen = tag["LastSeenTimestampUTC"]

            antenna = 1

            try:
                antenna = tag["AntennaID"]
            except Exception as e:
                logging.warning(
                    f"Exception getting antenna id: {e}",
                )

            last_seen = pendulum.now().in_timezone("America/Sao_Paulo")  # Damn.

            tag_obj: Tag = Tag(epc, last_seen, antenna=antenna)

            # TODO: tag_obj: Tag = Tag.from_timestamps(epc, last_seen, antenna=antenna)

            try:
                self.sender.send_tag(tag_obj)
            except Exception as e:
                logging.warn(f"Ignoring error at tag callback '{e}'...")

    def configure(self) -> bool:
        enabled_antennas = [0]

        import os

        host = os.getenv("READER_IP")

        if not host:
            host = "10.19.3.61"

        subtype = os.getenv("READER_SUBTYPE")

        modecap = 2

        if subtype == "R420":
            modecap = 0

        factory_args: dict[str, int | bool | list] = dict(
            duration=1,
            report_every_n_tags=1,  # TEST -> CONDITIONAL TAG REPORTS
            report_timeout_ms=0,
            antennas=enabled_antennas,
            tx_power=91,
            tari=0,
            session=2,
            tag_population=200,
            start_inventory=True,
            disconnect_when_done=False,
            reconnect=True,
            reconnect_retries=50,
            tag_filter_mask=[],
            tag_content_selector={
                "EnableROSpecID": False,
                "EnableSpecIndex": False,
                "EnableInventoryParameterSpecID": False,
                "EnableAntennaID": True,
                "EnableChannelIndex": False,
                "EnablePeakRSSI": False,
                "EnableFirstSeenTimestamp": True,
                "EnableLastSeenTimestamp": True,
                "EnableTagSeenCount": True,
                "EnableAccessSpecID": False,
                "C1G2EPCMemorySelector": {
                    "EnableCRC": False,
                    "EnablePCBits": False,
                },
            },
            frequencies={"HopTableId": 1, "ChannelList": [1], "Automatic": False},
            keepalive_interval=60000,
            impinj_extended_configuration=True,
            impinj_search_mode=modecap,
            mode_identifier=modecap,
            impinj_tag_content_selector=None,
        )

        # reader_clients = []

        if ":" in host:
            host, port = host.split(":", 1)
            port = int(port)
        else:
            port = 5084

        config: LLRPReaderConfig = LLRPReaderConfig(factory_args)
        reader: LLRPReaderClient = LLRPReaderClient(host, port, config)

        reader.add_disconnected_callback(self.finish_callback)
        reader.add_tag_report_callback(self.tag_report_callback)
        reader.add_state_callback(
            LLRPReaderState.STATE_INVENTORYING, self.inventory_start_callback
        )

        try:
            reader.connect()
        except Exception as e:
            logging.error(
                f"Error connecting to Reader: {e}",
            )
            return False

        self.reader: LLRPReaderClient = reader
        return True

    def read(self) -> None:
        try:
            self.reader.join()
        except (SystemExit, KeyboardInterrupt):
            logging.info("Stopping readers...")

    def start(self) -> Never:
        try:
            while not self.configure():
                logging.info("No Reader connected! reconnecting")
                sleep(2)

            logging.info("Reader connected! starting read")

            self.read()
            self.reader.disconnect()
            LLRPReaderClient.disconnect_all_readers()

            self.reader = None

        except Exception as e:
            logging.error(
                f"Reader error {e}",
            )
