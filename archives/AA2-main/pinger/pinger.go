package pinger

import (
	"sync/atomic"
	"time"

	//"aa2/lcdlogger"
	"github.com/prometheus-community/pro-bing"
)

func NewPinger(ip string, state *atomic.Bool, ping *atomic.Int64) {

	p, err := probing.NewPinger(ip)

	//p.SetPrivileged(true)

	if err != nil {

		return
	}

	p.Count = 0xFFFE
	p.Interval = 4 * time.Second

	p.OnSend = func(pkt *probing.Packet) {

		state.Store(false)
	}

	p.OnRecv = func(pkt *probing.Packet) {

		state.Store(true)

		if ping != nil {

			ping.Store(pkt.Rtt.Milliseconds())
		}
	}

	p.Run()
}
