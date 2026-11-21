package lcdlogger

import (
	"context"
	"log"
	"time"

	c "aa2/constant"

	"github.com/MyTempoesp/flick"
)

type SerialDisplay struct {
	Forth *flick.MyTempo_Forth

	Screen int

	actionButtonLHTime time.Time // Last Held Timestamp
	actionButtonHeld   bool
	action             Action

	altButtonLHTime time.Time
	altButtonHeld   bool
	altAction       Action
}

func NewSerialDisplay() (display SerialDisplay, err error) {

	f, err := flick.NewForth("/dev/ttyUSB0")

	if err != nil {

		log.Printf("Erro ao iniciar a comunicação com o arduino: %v\n", err)

		return
	}

	f.Start()

	f.Query("1 .")

	display.Forth = &f

	display.action = -1
	display.altAction = -1

	return
}

func (display *SerialDisplay) WaitKeyPress(d time.Duration) (hasKey bool) {

	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

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

				if res[0] == '-' {

					hasKey = true

					return
				}
			}
		}
	}
}

func (display *SerialDisplay) SwitchScreens() {

	// TODO: onrelease actions

	res, err := display.Forth.Send("bac @ .")

	if err != nil {

		return
	}

	defer display.Forth.Send("0 bac !")

	if res[0] == '-' { // onrelease

		lht := display.actionButtonLHTime

		if display.actionButtonHeld && time.Now().After(lht.Add(time.Millisecond*2300)) { // XXX: magic number

			display.action = Action(display.Screen)

			display.actionButtonLHTime = time.Now().Add(time.Hour * 24)

			display.actionButtonHeld = false // button released

			return
		}

		display.Screen++
		display.Screen %= SCREEN_COUNT
	}
}
