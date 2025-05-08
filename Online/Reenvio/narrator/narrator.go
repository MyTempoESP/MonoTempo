package narrator

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

func (n *Narrator) Close() {
	close(n.queue)
}

func (n *Narrator) Consume() {

	for {
		select {
		case s := <-n.queue:
			Say(s)
		default:
			return
		}
	}
}
