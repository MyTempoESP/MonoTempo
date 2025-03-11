# from collections import deque
from typing import Self

import pendulum

# from config.config import TMZ1
from helpers.functions import datetime_from_microseconds, to_microseconds

from .epc import EPC

MYTEMPO_TIME_FORMAT = "HH:mm:ss.SSS"  # Hour:Minute:Second.Millis


TMZ1 = "America/Sao_Paulo"


class Tag:
    def __init__(
        self,
        epc: EPC,
        capture_time: pendulum.DateTime,
        antenna: int = 1,  # antenna: Antenna
    ) -> None:
        self.epc: EPC = epc

        self.current_year: int = pendulum.now().year  # lol

        self.capture_time: pendulum.DateTime = pendulum.now()
        self.set_time_safely(capture_time)  # that's crazy

        self.antenna: int = antenna

    @classmethod
    def from_timestamps(
        make: Self, epc: EPC, capture_time: int, antenna: int = 1
    ) -> Self:
        return make(
            epc,
            datetime_from_microseconds(capture_time),
            antenna,
        )

    def formatted_time(self) -> str:
        return self.capture_time.format(MYTEMPO_TIME_FORMAT)

    def set_time_safely(self, time: pendulum.DateTime):
        if time.year == self.current_year:
            self.capture_time = time
        else:
            print(
                f"\033[31mBogus time at SET_TIME_SAFELY, {time}, setting to now()\033[0m"
            )
            self.capture_time = pendulum.now()

    def __str__(self) -> str:
        return f"{self.antenna:03}{self.epc.raw_epc[8:]}{self.capture_time.format(MYTEMPO_TIME_FORMAT)}"

    def __repr__(self) -> str:
        microseconds: int = int(to_microseconds(self.capture_time.float_timestamp))

        return f"{self.epc}:{microseconds}:{self.antenna}"
