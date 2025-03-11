package lcdlogger

const (
	ACTION_TAGS  = SCREEN_TAGS
	ACTION_WIFI  = SCREEN_WIFI
	ACTION_TIME  = SCREEN_TIME
	ACTION_USB   = SCREEN_USB
	ACTION_RESET = SCREEN_INFO_EQUIP
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
