package lcdlogger

import (
	"log"
	"time"

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

func (display *SerialDisplay) SwitchScreens() {

	// TODO: onrelease actions

	res, err := display.Forth.Send("bac @ .")

	if err != nil {

		return
	}

	defer display.Forth.Send("0 bac !")

	if res[0] == '-' { // onrelease

		display.actionButtonHeld = false // button released

		lht := display.actionButtonLHTime

		if time.Now().After(lht.Add(time.Millisecond * 1300)) { // XXX: magic number

			display.action = Action(display.Screen)

			display.actionButtonLHTime = time.Now().Add(time.Hour * 24)

			return
		}

		display.Screen++
		display.Screen %= SCREEN_COUNT
	}
}
