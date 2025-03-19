package lcdlogger

const (
	ACTION_TAGS          = SCREEN_TAGS
	ACTION_TIME          = SCREEN_TIME
	ACTION_USB           = SCREEN_USB
	ACTION_RESET         = SCREEN_INFO_EQUIP
	ACTION_UPLOAD        = SCREEN_UPLOAD
	ACTION_UPLOAD_BACKUP = SCREEN_UPLOAD_BACKUP
	ACTION_REBOOT        = SCREEN_ADDR
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
