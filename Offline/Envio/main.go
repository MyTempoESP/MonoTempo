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

func (envio *Envio) CreateBindings() (err error) {

	bindings :=
		[]rabbit.Binding{
			rabbit.NewBinding(
				ANTENNA_QUEUE,
				ANTENNA_ROUTE,
				ANTENNA_EXCHANGE,
				true, /* durable */
			),
		}

	err = envio.broker.BindQueues(bindings)

	return
}

func (envio *Envio) StartConsumers() (channel *amqp.Channel, err error) {

	channel, err = envio.broker.Channel()

	if err != nil {
		return
	}

	/*
		o Canal se mantém aberto para que se mantenha a leitura.
	*/
	//defer channel.Close()

	envio.channel = channel

	tags, err := ReceiveAntennas(channel)

	envio.Tags = tags

	return
}

func main() {
	for {
		var envio Envio

		err := envio.broker.Setup()

		if err != nil {

			log.Fatal(err)
		}

		err = envio.CreateBindings()

		if err != nil {

			log.Fatal(err)
		}

		channel, err := envio.StartConsumers()

		if err != nil {

			log.Fatal(err)
		}

		envio.DBManager.GroupSize(10)

		envio.DBManager.DatabaseRoot = "/var/monotempo-data"

		envio.DBManager.Init()

		go envio.Process()

		channelClosed := make(chan *amqp.Error)
		brokerClosed := make(chan *amqp.Error)

		channel.NotifyClose(channelClosed)
		envio.broker.NotifyClose(brokerClosed)

		select {
		case <-channelClosed:
			log.Println("Canal RabbitMQ fechado")
		case <-brokerClosed:
			log.Println("RabbitMQ encerrado")
			return
		}
	}
}
