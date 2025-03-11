package flick

import (
	"bytes"
	"fmt"
	"log"
	"sync"
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
	port         *serial.Port
	mu           sync.Mutex
	responseChan chan []byte
}

func NewForth(dev string, timeout time.Duration) (f Forth, err error) {

	conf := &serial.Config{
		Name:        dev,
		Baud:        115200,
		ReadTimeout: timeout,
	}

	f.port, err = serial.OpenPort(conf)

	if err != nil {

		log.Fatalf("Failed to open serial port: %v", err)
	}

	f.responseChan = make(chan []byte)

	return
}

func (f *Forth) Stop() {

	f.port.Close()
	close(f.responseChan)
}

func (f *Forth) Start() {

	go func() {

		buf := make([]byte, 128)

		for {
			n, err := f.port.Read(buf)

			if err != nil {

				f.responseChan <- []byte("(timeout!)")

				continue
			}

			if n > 0 {

				f.responseChan <- buf[:n]
			}
		}
	}()
}

func (f *Forth) Send(input string) (response []byte, err error) {

	f.mu.Lock()
	defer f.mu.Unlock()

	_, err = f.port.Write([]byte(input + "\n"))

	if err != nil {

		log.Printf("Failed to send data: %v", err)

		return
	}

	response = <-f.responseChan

	return
}

func (f *Forth) Query(input string) (multilineResponse []byte, err error) {

	f.mu.Lock()
	defer f.mu.Unlock()

	_, err = f.port.Write([]byte(input + "\n"))

	if err != nil {

		log.Printf("Failed to send data: %v", err)

		return
	}

	fmt.Printf("Sent: %s\n", input)

	for {
		response := <-f.responseChan
		multilineResponse = append(multilineResponse, response...)

		if bytes.Contains(response, []byte{'o', 'k'}) {

			break
		}
	}

	return
}
