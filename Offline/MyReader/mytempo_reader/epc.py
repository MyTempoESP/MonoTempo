import re
from typing import Self


class InvalidEPCError(Exception):
    pass


class EPC:
    def __init__(self, epc: str):
        split_epc: list[str] = re.findall("[0-9]{4}", epc)
        if len(split_epc) < 6 or "abcdefABCDEF" in epc:
            raise InvalidEPCError(
                "Wrong format for epc (check if the tag has invalid hexadecimal values)"
            )

        self.raw_epc = epc  # for Refined data format
        self.epc: str = "-".join(split_epc)

        self.id: int = int(split_epc[-1])  # only last 4 digits matter for id

    @classmethod
    def from_tag_data(make: Self, tag_data: dict) -> Self:
        return make(tag_data["EPC"].decode("utf-8"))

    def __str__(self) -> str:
        return self.epc

    def __repr__(self) -> str:
        return f"EPC({self.epc}, {hash(self.id)})"

    def __hash__(self) -> hash:
        return hash(self.id)

    def __eq__(self, other: any) -> bool:
        if isinstance(other, EPC):
            return self.id == other.id

        return False
