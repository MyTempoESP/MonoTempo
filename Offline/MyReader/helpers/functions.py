from random import randrange

import pendulum


def make_sess(session_file_path: str) -> str:
    return f"{session_file_path}mytempo-Sess-{randrange(0,999):03}.txt"


def to_seconds(microseconds: int):
    return microseconds / 1e6


def to_microseconds(seconds: int):
    return seconds * 1e6


def datetime_from_microseconds(microseconds: int) -> pendulum.DateTime:
    date = pendulum.from_timestamp(to_seconds(microseconds), tz="UTC").in_timezone(
        "America/Sao_Paulo"
    )
    # return date.format(MYTEMPO_TIME_FORMAT)
    return date
