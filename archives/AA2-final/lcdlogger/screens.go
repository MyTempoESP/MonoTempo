package lcdlogger

import (
	"fmt"
	"log"
	"time"

	c "aa2/constant"
)

const (
	SCREEN_TAGS = iota
	SCREEN_IP
	SCREEN_WIFI
	SCREEN_ANTENNAS
	SCREEN_TIME
	SCREEN_USB
	SCREEN_SIST
)

type IPOctets [4]int

func (display *SerialDisplay) DrawScreen(code string) {

	display.Forth.Send("0 TME")
	display.Forth.Send(code)
	display.Forth.Send("1 TME")
}

func (display *SerialDisplay) ScreenTags(tags, tagsUnicas ForthNumber) {

	display.DrawScreen(
		fmt.Sprintf(
			"%d %d %d %d S-1",

			tags.Value, tags.Magnitude,
			tagsUnicas.Value, tagsUnicas.Magnitude,
		),
	)
}

func (display *SerialDisplay) ScreenAddr(ip IPOctets, leitorOk int) {

	log.Println(ip)

	display.DrawScreen(
		fmt.Sprintf(
			"%d %d %d %d %d S-2",

			leitorOk,
			ip[0], ip[1], ip[2], ip[3],
		),
	)
}

func (display *SerialDisplay) ScreenWifi(wifiVerif int, wifiPing int64) {

	display.DrawScreen(
		fmt.Sprintf(
			"%d %d 6 S-3",

			wifiPing,
			wifiVerif,
		),
	)
}

func (display *SerialDisplay) ScreenStat(a1, a2, a3, a4 ForthNumber) {

	display.DrawScreen(
		fmt.Sprintf(
			"%d %d"+ // A4 Val+Mag
				" %d %d"+ // A3 Val+Mag
				" %d %d"+ // A2 Val+Mag
				" %d %d S-4", // A1 Val+Mag then display

			a4.Magnitude, a4.Value,
			a3.Magnitude, a3.Value,
			a2.Magnitude, a2.Value,
			a1.Magnitude, a1.Value,
		),
	)
}

func (display *SerialDisplay) ScreenTime() {

	now := time.Now().In(c.ProgramTimezone)
	y, m, d := now.Date()
	//log.Println("now", now)
	//hour min sec day mon year

	display.DrawScreen(
		fmt.Sprintf(
			"%d %d %d %d %d %d S-5",

			now.Hour(), now.Minute(), now.Second(),
			d, m, y,
		),
	)
}

func (display *SerialDisplay) ScreenUSB(devVerif int) {

	display.DrawScreen(
		fmt.Sprintf(
			"%d S-6",
			devVerif,
		),
	)
}

func (display *SerialDisplay) ScreenInfoEquip() {

	display.DrawScreen(
		fmt.Sprintf(
			"$%s S-7",

			//c.Reader, DAMN
			c.Version,
		),
	)
}

func (display *SerialDisplay) ScreenConfirma() {

	display.DrawScreen(
		fmt.Sprintf(
			"S-8",
		),
	)
}

func (display *SerialDisplay) ScreenProgress() {

	display.DrawScreen(
		fmt.Sprintf(
			"",
		),
	)
}

func (display *SerialDisplay) ScreenErr() {

	display.DrawScreen(
		fmt.Sprintf(
			"S-9",
		),
	)
}
