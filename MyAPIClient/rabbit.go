package main

import (
        amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	host string,
	port string,

	user string,
	pass string,
}


func bodyFrom(args []string) string {
	var s string
        if (len(args) < 3) || os.Args[2] == "" {
                s = "hello"
        } else {
                s = strings.Join(args[2:], " ")
        }
        return s
}


func severityFrom(args []string) string {
        var s string
        if (len(args) < 2) || os.Args[1] == "" {
                s = "anonymous.info"
        } else {
                s = os.Args[1]
        }
        return s
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


func (rabbit *Rabbit) setupAuth() {
	user, pass, err := getRabbitAuth()
	failOnError(err, "Failed to get RabbitMQ credentials")

	rabbit.user = user
	rabbit.pass = pass
}

func (rabbit *Rabbit) setupServer() {
	host, port, err := getRabbitServer()
	failOnError(err, "Failed to get RabbitMQ host")

	rabbit.host = host
	rabbit.port = port
}

func (rabbit *Rabbit) Url() {
	auth := fmt.Sprintf("%s:%s", rabbit.user, rabbit.pass)
	url := fmt.Sprintf("amqp://%s@%s:%s/", auth, rabbit.host, rabbit.port)
}

func (rabbit *Rabbit) connect() {
	url = rabbit.Url()


	rabbit.port = 
}

func (rabbit *Rabbit) Setup() {
	rabbit.setupServer()
	rabbit.setupAuth()

	rabbit.connect()

	conn, err := amqp.Dial(url)
        failOnError(err, "Failed to connect to RabbitMQ")
        defer conn.Close()

        channel, err := conn.Channel()
        failOnError(err, "Failed to open a channel")
        defer channel.Close()

        err = channel.ExchangeDeclare(
                "logs_topic", // name
                "topic",      // type
                true,         // durable
                false,        // auto-deleted
                false,        // internal
                false,        // no-wait
                nil,          // arguments
        )

        failOnError(err, "Failed to declare an exchange")

        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        body := bodyFrom(os.Args)
        err = channel.PublishWithContext(ctx,
                "logs_topic",          // exchange
                severityFrom(os.Args), // routing key
                false, // mandatory
                false, // immediate
                amqp.Publishing{
                        ContentType: "text/plain",
                        Body:        []byte(body),
                })
        failOnError(err, "Failed to publish a message")

        log.Printf(" [x] Sent %s", body)
}
