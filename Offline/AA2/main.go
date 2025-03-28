package main

/*
TODO: logging module for log levels and passing
log objects to functions so they don't really bother us
*/

import (
	"aa2/constant"
	"log"
	"os/exec"

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

func filhoDaPutaVaiSeFuderArrombado(e error) {
	cmd := exec.Command("sh", "-c", "echo 'fatal' > /var/monotempo-data/sig-upload-data")
	err := cmd.Run()
	log.Println(err, e)
	select {}
}

func main() {
	if constant.VersionNum == "" {
		log.Println(constant.VersionNum)
		constant.VersionNum = "0000"
	}

	for {
		var a Ay

		err := a.broker.Setup()

		if err != nil {
			filhoDaPutaVaiSeFuderArrombado(err)
		}

		err = a.CreateBindings()

		if err != nil {
			filhoDaPutaVaiSeFuderArrombado(err)
		}

		channel, err := a.StartConsumers()

		if err != nil {
			filhoDaPutaVaiSeFuderArrombado(err)
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
