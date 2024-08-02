package main

// errors

import (
	rabbit "github.com/mytempoesp/rabbit"
)

func main() {
	var r rabbit.Rabbit

	r.Setup()
	//r.NewTopic("api_exchange") unreleased
	r.SendMessage("Oie", 10)
}
