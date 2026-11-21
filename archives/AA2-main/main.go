package main

/*
TODO: logging module for log levels and passing
log objects to functions so they don't really bother us
*/

import (
	"log"

	rabbit "github.com/mytempoesp/rabbit"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (a *Ay) CreateBindings() (err error) {

	bindings :=
		[]rabbit.Binding{
			rabbit.NewBinding(
				ANTENNA_QUEUE,
				ANTENNA_ROUTE,
				ANTENNA_EXCHANGE,
				true, /* durable */
			),
		}

	err = a.broker.BindQueues(bindings)

	return
}

func (a *Ay) StartConsumers() (channel *amqp.Channel, err error) {

	channel, err = a.broker.Channel()

	if err != nil {
		return
	}

	a.channel = channel

	tags, err := ReceiveAntennas(channel)

	a.Tags = tags

	return
}

func main() {
	for {
		var a Ay

		err := a.broker.Setup()

		if err != nil {
			log.Fatal(err)
		}

		err = a.CreateBindings()

		if err != nil {
			log.Fatal(err)
		}

		channel, err := a.StartConsumers()

		if err != nil {
			log.Fatal(err)
		}

		go a.Process()

		channelClosed := make(chan *amqp.Error)
		brokerClosed := make(chan *amqp.Error)

		channel.NotifyClose(channelClosed)
		a.broker.NotifyClose(brokerClosed)

		select {
		case <-channelClosed:
			log.Println("Canal RabbitMQ fechado")
		case <-brokerClosed:
			log.Println("RabbitMQ encerrado")
			return
		}
	}
}
