package lcdlogger

import (
	"aa2/flick"
	"log"
	"time"

	c "aa2/constant"
)

type SerialDisplay struct {
	Forth *flick.Forth

	Screen int
	action Action
}

func NewSerialDisplay() (display SerialDisplay, err error) {

	f, err := flick.NewForth(c.ReaderPath, 2*time.Second)

	if err != nil {

		log.Printf("Erro ao iniciar a comunicação com o arduino: %v\n", err)

		return
	}

	f.Send("1 .")

	display.Forth = &f
	display.action = -1

	return
}

/*
func (display *SerialDisplay) WaitKeyPress(d time.Duration) (hasKey bool) {

	ctx, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()

	for {
		select {
		case <-ctx.Done():

			hasKey = false

			goto done
		default:
			{
				res, err := display.Forth.Send(c.FORTH_BTN_PRESSED)

				if err != nil {

					continue
				}

				if res[0] == '1' {

					hasKey = true

					goto done
				}
			}
		}
	}

done:
	display.Forth.Send("0 bac ! 0 ba2 !")

	return
}

func (display *SerialDisplay) SwitchScreens() {

	// TODO: onrelease actions

	<-time.After(100 * time.Millisecond)

	res, err := display.Forth.Send("ba@ b2@")

	if err != nil {

		return
	}

	res = bytes.TrimSuffix(res, []byte{' ', 'o', 'k', '\n'})

	m, err := regexp.Match("^[0-1] [0-1]$", res)

	if err != nil || !m {

		log.Printf("Unexpected output, expected '^[0-1] [0-1]$', got '%s'\n", res)

		return
	}

	if res[0] == '1' {

		log.Println("Switching screen!")

		display.Screen++
		display.Screen %= SCREEN_COUNT

		display.Forth.Send("0 bac !")
	}

	if res[2] == '1' {

		log.Println("Confirm!")

		display.action = Action(display.Screen)

		display.Forth.Send("0 ba2 !")
	}
}
*/
