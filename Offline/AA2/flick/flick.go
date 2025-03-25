package flick

import (
	"fmt"
	"time"

	"github.com/MyTempoESP/serial"
)

// LABELS
const (
	PORTAL = iota
	UNICAS
	REGIST
	COMUNICANDO
	LEITOR
	LTE4G
	WIFI
	IP
	LOCAL
	PROVA
	PING

	LABELS_COUNT
)

// VALUES
const (
	WEB = iota
	CONECTAD
	DESLIGAD
	AUTOMATIC
	OK
	X

	VALUES_COUNT
)

type Forth struct {
	port *serial.Port
}

func NewForth(dev string, timeout time.Duration) (f Forth, err error) {

	conf := &serial.Config{
		Name:        dev,
		Baud:        115200,
		ReadTimeout: timeout,
	}

	f.port, err = serial.OpenPort(conf)

	return
}

func (f *Forth) Stop() {

	f.port.Close()
}

func (f *Forth) Send(input string) (response []byte, err error) {

	_, err = f.port.Write(fmt.Appendf([]byte{}, "%c%s%c\n", 0x3C, input, 0x3E))

	<-time.After(100 * time.Millisecond)
	f.port.Read(response)

	return
}
