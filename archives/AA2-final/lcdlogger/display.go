package lcdlogger

import (
	"bytes"
	"context"
	"log"
	"strconv"
	"time"

	c "aa2/constant"

	"aa2/flick"
)

const (
	NOT_READY = -1
)

type SerialDisplay struct {
	Forth *flick.Forth

	Screen int
	action Action
}

func NewSerialDisplay() (display SerialDisplay, err error) {

	f, err := flick.NewForth("/dev/ttyACM0", 500*time.Millisecond)

	if err != nil {

		log.Printf("Erro ao iniciar a comunicação com o arduino: %v\n", err)

		return
	}

	f.Start()

	f.Query("1 .")

	display.Forth = &f

	display.action = -1

	return
}

func (display *SerialDisplay) WaitKeyPress(d time.Duration) (hasKey bool) {

	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	defer display.Forth.Send("AC@ .")

	for {
		select {
		case <-ctx.Done():

			hasKey = false

			return
		default:
			{
				res, err := display.Forth.Send(c.FORTH_BTN_PRESSED)

				if err != nil {

					continue
				}

				if res[0] != '0' {

					hasKey = true

					return
				}
			}
		}
	}
}

func (display *SerialDisplay) SwitchScreens() (n int) {

	// TODO: onrelease actions

	res, err := display.Forth.Query("2 BA@ .") // fetch action

	if err != nil {

		return
	}

	res = bytes.TrimSuffix(res, []byte{' ', 'o', 'k', '\n'})

	log.Printf("'%s'", res)

	n64, err := strconv.ParseInt(string(res), 16, 8)

	if err != nil {
		return
	}

	n = int(n64) & 0x0F

	// display.action = Action(n)

	return
}
