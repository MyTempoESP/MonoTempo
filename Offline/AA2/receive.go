package main

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func AMQP2Tag(body []byte) (t tag, err error) {

	err = json.Unmarshal(body, &t)

	return
}

func NewExclusiveConsumer(queue string, name string, channel *amqp.Channel) (recv <-chan tag, err error) {

	consume, err := channel.Consume(
		queue, // queue
		name,  // consumer
		false, // autoAck
		true,  // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)

	if err != nil {
		log.Printf("basic.consume: %v", err)
		return
	}

	receiver := make(chan tag)

	go func() {
		for delivery := range consume {
			err := delivery.Ack(false)

			if err != nil {
				log.Printf("erro de ack: %+v", err)

				continue
			}

			//log.Printf("Got data! %s\n", string(delivery.Body))

			t, err := AMQP2Tag(delivery.Body)

			if err != nil {
				log.Printf("Não foi possível interpretar a tag em JSON: %+v", err)

				continue
			}

			receiver <- t
		}

		close(receiver)
	}()

	recv = receiver

	return
}

func ReceiveAntennas(channel *amqp.Channel) (<-chan tag, error) {

	return NewExclusiveConsumer(
		ANTENNA_QUEUE, ANTENNA_CONSUMER_NAME, channel,
	)
}
