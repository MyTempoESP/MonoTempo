package lcdlogger

import (
	"time"
)

const (
	ALT_ACTION_RESTART = iota
)

func (display *SerialDisplay) AltAction() (action Action, hasAction bool) {

	action = display.altAction

	if action >= 0 {

		hasAction = true
		display.altAction = -1
	}

	return
}

func (display *SerialDisplay) hasAltAction() bool {

	return display.altAction >= 0
}

func (display *SerialDisplay) HandleAltActionButton() {

	if display.altButtonHeld {

		return
	}

	if display.hasAltAction() {

		return
	}

	res, err := display.Forth.Send("bs2 @ .") // ButtonSTate

	if err != nil {

		return
	}

	if res[0] == '-' {

		display.altButtonLHTime = time.Now()
		display.altButtonHeld = true
	}
}

// Alternative action button
func (display *SerialDisplay) ProcessAltAction() {

	// TODO: onrelease actions

	res, err := display.Forth.Send("ba2 @ .")

	if err != nil {

		return
	}

	defer display.Forth.Send("0 ba2 !")

	if res[0] == '-' { // onrelease

		lht := display.altButtonLHTime

		if time.Now().After(lht.Add(time.Millisecond * 1300)) { // XXX: magic number

			display.altAction = Action(ALT_ACTION_RESTART)

			display.altButtonLHTime = time.Now().Add(time.Hour * 24)

			return
		}
	}
}
