package narrator

import (
	"time"
)

type Narrator struct {
	Enabled bool

	queue chan string
}

func New() (n Narrator) {

	n.Enabled = true
	n.queue = make(chan string, 10)

	return
}

func (n *Narrator) SayString(s string) {
	n.queue <- s
}

func (n *Narrator) Watch() {

	for s := range n.queue {

		Say(s)

		<-time.After(5 * time.Second)
	}
}
