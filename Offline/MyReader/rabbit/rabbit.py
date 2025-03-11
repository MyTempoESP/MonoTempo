import json
import logging
import time
from typing import override

import pika
from myconfig import RABBIT_HOST, RABBIT_PASS, RABBIT_USER
from mytempo_reader.tag import Tag

logging.basicConfig(level=logging.INFO)


class MQClient:

    def __init__(self, exchange="", host="localhost", heartbeat_interval=1800):
        self.host = RABBIT_HOST
        self.heartbeat_interval = heartbeat_interval

        self.credentials = pika.PlainCredentials(RABBIT_USER, RABBIT_PASS)

        self.connection = None
        self.channel = None

    def connect(self):
        while True:
            try:
                self.connection = pika.BlockingConnection(
                    pika.ConnectionParameters(
                        self.host,
                        heartbeat=self.heartbeat_interval,
                        credentials=self.credentials,
                    )
                )

                return

            except pika.exceptions.AMQPConnectionError as e:
                logging.error(f"Connection failed: {e}. Retrying in 1 seconds...")
                time.sleep(1)  # Wait before retrying

    def new_channel(self):
        while True:
            try:
                self.channel = self.connection.channel()
                break

            except Exception as e:
                logging.error(f"Channel creation failed: {e}.")
                self.connect()
                time.sleep(1)

        logging.info("Connected to RabbitMQ")

    def declare_exchanges(self): ...

    def declare_exchange(self, exchangeName, isDurable):
        self.channel.exchange_declare(
            exchange=exchangeName, durable=isDurable, exchange_type="topic"
        )

    def declare_queues(self): ...

    def declare_queue(self, queueName, isDurable):
        self.channel.queue_declare(queue=queueName, durable=isDurable)

    def bind(self, queueName, exchangeName, routing_key):
        self.channel.queue_bind(
            exchange=exchangeName,
            queue=queueName,
            routing_key=routing_key,
        )

    def start(self): ...

    def setup(self):
        self.connect()
        self.new_channel()

        self.declare_exchanges()
        self.declare_queues()

    def publish(self, message: str, routing_key, exchange):
        try:
            if not self.channel or self.channel.is_closed:
                self.new_channel()  # Reconnect if the channel is closed

            self.channel.basic_publish(
                exchange=exchange,
                routing_key=routing_key,
                body=message.encode(),
                mandatory=True,
            )

            # logging.info(f"Sent message: {message}")

        except pika.exceptions.UnroutableError:
            logging.warning(f"Message {message} could not be routed.")

        except pika.exceptions.AMQPChannelError as e:
            logging.error(f"Channel error: {e}. Reconnecting...")
            self.connect()

        except Exception as e:
            logging.error(f"Unknown error in rabbitmq {e}. Reconnecting")
            self.connect()

    def subscribe(self, callback, queueName):
        if not self.channel or self.channel.is_closed:
            self.new_channel()

        try:
            self.channel.basic_consume(
                queue=queueName,
                auto_ack=True,
                on_message_callback=callback,
            )

        except pika.exceptions.AMQPChannelError as e:
            logging.error(f"Channel error: {e}. Reconnecting...")
            self.connect()
            self.new_channel()

    def __enter__(self):
        self.setup()

        return self

    def __exit__(self, *exc):
        if self.connection:
            self.connection.close()
            logging.info("Connection closed.")


class RFIDPublisher(MQClient):

    @override
    def declare_exchanges(self):
        super().declare_exchanges()

        self.declare_exchange("antenna_exchange", True)


class TagSender(RFIDPublisher):

    def send_tag(self, tag: Tag):
        logging.info(f"Sending: {repr(tag)}")
        routing_key = f"antenna.{tag.antenna}"

        tag_json = {
            "antena": tag.antenna,
            "refinado_mytempo": str(tag),
            "tempo_formatado": tag.formatted_time(),
            "tempo_formato": "15:04:05.000",
            "epc": tag.epc.id,
        }

        self.publish(json.dumps(tag_json), routing_key, "antenna_exchange")

    @override
    def declare_queues(self):
        super().declare_queues()

        for n in range(1, 5):  # 1 to 4 inclusive
            queue = f"antenna_{n}_queue"
            routing_key = f"antenna.{n}"

            self.declare_queue(queue, True)
            self.bind(queue, "antenna_exchange", routing_key)
