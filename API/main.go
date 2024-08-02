package main

// errors

import (
	//"log
	rabbit "github.com/mytempoesp/rabbit"
)

/*
our own fail function will be required
in v0.0.2

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
*/

func main() {
	var r rabbit.Rabbit

	r.Setup()
	//r.NewTopic("api_exchange")
	r.SendMessage("Oie", 10)
}
