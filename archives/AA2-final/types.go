package main

import (
	rabbit "github.com/mytempoesp/rabbit"
	amqp "github.com/rabbitmq/amqp091-go"
)

type tag struct {
	Epc             int    `json:"epc"`
	Antena          int    `json:"antena"` /* min: 1, max: 4 */
	Staff           int    `json:"staff"`
	TempoFormato    string `json:"tempo_formato"`
	TempoFormatado  string `json:"tempo_formatado"`
	FormatoRefinado string `json:"refinado_mytempo"`
}

type Ay struct {
	broker  rabbit.Rabbit
	channel *amqp.Channel

	Tags <-chan tag
}
