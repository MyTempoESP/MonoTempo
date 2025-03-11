from os import getenv


class MissingEnvError(Exception): ...


RABBIT_HOST = getenv("RABBITMQ_HOST")
RABBIT_USER = getenv("RABBITMQ_USER")
RABBIT_PASS = getenv("RABBITMQ_PASS")

if not any((RABBIT_HOST, RABBIT_USER, RABBIT_PASS)):
    raise MissingEnvError("RABBIT envs not properly setup")
