package narrator

import (
	"strconv"
	"time"
)

type Narrator struct {
	Enabled bool

	queue chan int
}

func New() (n Narrator) {

	n.Enabled = true
	n.queue = make(chan int, 200)

	return
}

func (n *Narrator) SayNum(id int) {
	n.queue <- id
}

func (n *Narrator) Watch() {

	for id := range n.queue {

		Say(strconv.Itoa(id))

		<-time.After(2 * time.Second)
	}
}
