package main

import (
	"fmt"
	"log"
)

var ErrNoRabbitKey = fmt.Errorf("RABBITMQ_PASS env not found")
var ErrNoRabbitUser = fmt.Errorf("RABBITMQ_USER env not found")
var ErrNoRabbitHost = fmt.Errorf("RABBITMQ_HOST env not found")


func failOnError(err error, msg string) {
        if err != nil {
                log.Panicf("%s: %s", msg, err)
        }
}


