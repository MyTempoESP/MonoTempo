package lcdlogger

import (
	"time"
)

const (
	ACTION_TAGS = SCREEN_TAGS
	ACTION_WIFI = SCREEN_WIFI
	ACTION_TIME = SCREEN_TIME
	ACTION_USB  = SCREEN_USB
)

type Action int

func (display *SerialDisplay) Action() (action Action, hasAction bool) {

	action = display.action

	if action >= 0 {

		hasAction = true
		display.action = -1
	}

	return
}

func (display *SerialDisplay) hasAction() bool {

	return display.action >= 0
}

func (display *SerialDisplay) HandleActionButton() {

	if display.actionButtonHeld {

		return
	}

	if display.hasAction() {

		return
	}

	res, err := display.Forth.Send("bst @ .") // ButtonSTate

	if err != nil {

		return
	}

	if res[0] == '-' {

		display.actionButtonLHTime = time.Now()
		display.actionButtonHeld = true
	}
}
