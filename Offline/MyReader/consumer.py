import pika

import myconfig

credentials = pika.PlainCredentials("rodrigo", myconfig.RABBIT_TOKEN)


def on_message(ch, method, properties, body: bytes):
    print(f"Inventory Management received: {body.decode()}")
    ch.basic_ack(delivery_tag=method.delivery_tag)


connection = pika.BlockingConnection(
    pika.ConnectionParameters("localhost", credentials=credentials)
)
channel = connection.channel()

channel.exchange_declare(exchange="antenna_exchange", exchange_type="topic")

result = channel.queue_declare(queue="", exclusive=True)
queue_name = result.method.queue

channel.queue_bind(exchange="antenna_exchange", queue=queue_name, routing_key="antenna.#")

channel.basic_consume(queue=queue_name, on_message_callback=on_message)

print("Waiting for tags. To exit press CTRL+C")
channel.start_consuming()
