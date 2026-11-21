package lcdlogger

const (
	ACTION_TAGS = iota
	ACTION_IP
	ACTION_WIFI
	ACTION_ANTENNAS
	ACTION_TIME
	ACTION_USB
	ACTION_RESET
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
