package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	backoff "github.com/cenkalti/backoff"
	amqp "github.com/rabbitmq/amqp091-go"
)

/* TODO's:
- Cover the whole code with logging
*/

type Rabbit /* RabbitMQ connection */ struct {
	host string
	port string

	user string
	pass string

	exchange string /* current exchange name */

	routing_key string /*select Queue to which msg will b sent*/

	conn *amqp.Connection /* current open connection (Connection pointer) */
}

func getRabbitAuth() (user string, pass string, err error) {

	pass, err = os.Getenv("RABBITMQ_PASS"), nil

	if pass == "" {
		err = ErrNoRabbitKey
		return
	}

	user = os.Getenv("RABBITMQ_USER")

	if user == "" {
		err = ErrNoRabbitUser

		/* fallthrough */
	}

	return
}

func getRabbitServer() (host string, port string, err error) {

	err = nil

	host = os.Getenv("RABBITMQ_HOST")
	port = os.Getenv("RABBITMQ_PORT")

	if host == "" {
		err = ErrNoRabbitHost
		return
	}

	// FIXME: issue warning or something

	if port == "" {
		port = "5672"

		/* fallthrough */
	}

	return
}

func (rabbit *Rabbit) setupAuth() (err error) {

	var user, pass string

	user, pass, err = getRabbitAuth()

	if err != nil {
		return err
	}

	/* expose state */
	rabbit.user = user
	rabbit.pass = pass

	return nil
}

func (rabbit *Rabbit) setupServer() (err error) {

	var host, port string

	host, port, err = getRabbitServer()

	if err != nil {
		return err
	}

	/* expose state */
	rabbit.host = host
	rabbit.port = port

	return nil
}

func (rabbit *Rabbit) Url() (url string) {
	/*
		This function uses:

			rabbit.user
			rabbit.pass

			rabbit.host
			rabbit.port

		to generate an amqp url.
	*/

	var auth string

	/*
		setup a url in the form:

		amqp://<user>:<password>@<host>:<port>
	*/

	auth = fmt.Sprintf("%s:%s", rabbit.user, rabbit.pass)
	url = fmt.Sprintf("amqp://%s@%s:%s/", auth, rabbit.host, rabbit.port)

	return
}

func (rabbit *Rabbit) LogUrl() (url string) {
	/*
		Log- variation of Url function,
		this function doesn't expose the
		user password to the whole world
		:)
	*/

	/*
		This function uses:

			rabbit.user
			-rabbit.pass-

			rabbit.host
			rabbit.port

		to generate an amqp url. (Omitting password)
	*/

	var auth string

	/*
		setup a url in the form:

		amqp://<user>:******@<host>:<port>
	*/

	auth = fmt.Sprintf("%s:******", rabbit.user)
	url = fmt.Sprintf("amqp://%s@%s:%s/", auth, rabbit.host, rabbit.port)

	return
}

func (rabbit *Rabbit) Channel() (channel *amqp.Channel, err error) {

	log.Println("Opening channel")

	channel, err = rabbit.conn.Channel()

	if err != nil {
		return
	}

	return
}

func (rabbit *Rabbit) ExchangeDeclare(exchangeName string, exchangeType string, isDurable bool) (err error) {

	/* NOTE:
	This function opens a Channel.
	*/

	var channel *amqp.Channel

	/*

		Declare a RabbitMQ exchange.
		returns an error if the exchange
		could not be created, otherwise the
		exchange is created and the exchange is
		set accordingly.

		sets the 'rabbit.exchange' field to the
		created exchange name.

	*/

	log.Printf("Declaring exchange %s", exchangeName)

	channel, err = rabbit.Channel()

	if err != nil {
		return
	}

	defer channel.Close()

	err = channel.ExchangeDeclare(
		exchangeName,
		exchangeType,
		isDurable,
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return /* don't set the rabbit exchange */
	}

	rabbit.exchange = exchangeName

	return
}

func (rabbit *Rabbit) QueueDeclare(queueName string, isDurable bool) (err error) {

	/* NOTE:
	This function opens a Channel.
	*/

	var channel *amqp.Channel

	/*
		Declare a RabbitMQ Queue.

		The resulting queue is defined but not returned,
		it should be accessed using a 'routing_key' which
		matches the name of the queue.

		sets the 'rabbit.routing_key' field to the
		created queue name.

		the argument isDurable decides if the queue messages
		will persist after a broker restart

		returns an error status in case the Queue cannot be created.
	*/

	log.Printf("Declaring queue %s", queueName)

	channel, err = rabbit.Channel()

	if err != nil {
		return
	}

	defer channel.Close()

	_, err = channel.QueueDeclare(
		queueName,
		isDurable,
		false, // delete when unused (auto-deleted)
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	rabbit.routing_key = queueName

	return
}

func (rabbit *Rabbit) SendMessage(msg string, timeout_s int) (err error) {

	var channel *amqp.Channel

	log.Println("Trying to send message through RabbitMQ")

	timeout_duration := time.Duration(timeout_s)

	ctx, cancel := context.WithTimeout(context.Background(), timeout_duration*time.Second)
	defer cancel()

	channel, err = rabbit.Channel()

	if err != nil {
		return
	}

	defer channel.Close()

	log.Printf("Publishing message to exchange %s, queue(s): %s", rabbit.exchange, rabbit.routing_key)

	err = channel.PublishWithContext(ctx,
		rabbit.exchange,
		rabbit.routing_key,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain", /* "application/json", */
			Body:        []byte(msg),

			/*	NOTE:
				for Persistant messaging, uncomment
				the following line:
			*/
			//DeliveryMode: 2, /* default: 0 (transient) */
		})

	return
}

func (rabbit *Rabbit) Connect() (err error) {

	var url string
	var conn *amqp.Connection /* Connection pointer */

	url = rabbit.Url()

	log.Printf("Initiating connection attempt to '%s'", rabbit.LogUrl())

	err = backoff.Retry(
		func() error {
			log.Println("Trying connection to RabbitMQ")

			conn, err = amqp.Dial(url)

			if err != nil {
				log.Printf("Failed to connect to RabbitMQ: %v.\nRetrying...", err)
			}

			return err
		},

		backoff.NewExponentialBackOff(),
	)

	if err != nil {
		err = fmt.Errorf("Couldn't connect to RabbitMQ: %s", err)

		return
	}

	rabbit.conn = conn

	log.Printf("Succesfully connected to RabbitMQ at '%s'", rabbit.LogUrl())

	return
}

func (rabbit *Rabbit) Setup() {
	err := rabbit.setupServer()
	failOnError(err, "Failed to setup server")

	err = rabbit.setupAuth()
	failOnError(err, "Failed to setup authentication")

	err = rabbit.Connect()
	failOnError(err, "Failed to connect to RabbitMQ")

	err = rabbit.ExchangeDeclare(
		"api_exchange", // name
		"topic",        // type
		true,           // durable?
	)
	failOnError(err, "Failed to declare an exchange")

	err = rabbit.QueueDeclare(
		"api_data", // name
		true,       // durable?
	)
	failOnError(err, "Failed to declare a queue")

	//failOnError(err, "Failed to declare an exchange")
	//failOnError(err, "Failed to publish a message")

	//log.Printf(" [x] Sent %s", body)
}
