package narrator

import (
	"time"
)

type nothing struct{}

type Narrator struct {
	Enabled bool

	said  map[string]nothing
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

func (n *Narrator) Close() {
	close(n.queue)
}

func (n *Narrator) Consume() {

	for {
		select {
		case s := <-n.queue:

			_, exists := n.said[s]

			// do not repeat messages
			if exists {
				continue
			}

			n.said[s] = nothing{}

			Say(s)
		default:
			return
		}

		<-time.After(3 * time.Second)
	}
}
