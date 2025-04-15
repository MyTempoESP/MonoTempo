package pinger

import (
	"sync/atomic"
	"time"

	//"aa2/lcdlogger"
	probing "github.com/prometheus-community/pro-bing"
)

func NewPinger(ip string, state *atomic.Bool, ping *atomic.Int64) (p *probing.Pinger) {

	p, err := probing.NewPinger(ip)

	//p.SetPrivileged(true)

	if err != nil {

		return
	}

	p.Interval = 4 * time.Second

	p.OnRecv = func(pkt *probing.Packet) {

		if state != nil {
			state.Store(true)
		}

		if ping != nil {

			ping.Store(pkt.Rtt.Milliseconds())
		}
	}

	p.OnSendError = func(pkt *probing.Packet, err error) {
		if state != nil {
			state.Store(false)
		}
	}

	return
}
