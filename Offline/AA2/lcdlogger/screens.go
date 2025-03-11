package lcdlogger

import (
	"fmt"
	"log"
	"time"

	"github.com/MyTempoesp/flick"

	c "aa2/constant"
)

const (
	SCREEN_TAGS = iota
	SCREEN_ADDR
	SCREEN_WIFI
	SCREEN_STAT
	SCREEN_TIME
	SCREEN_USB
	SCREEN_INFO_EQUIP

	SCREEN_COUNT
)

const ( /* Labels Extras */
	LABEL_PROGRESSO = 13 + iota
	LABEL_ERRO
	LABEL_ERRO2

	LABEL_RFID
	LABEL_SERIE
	LABEL_SIST

	LABEL_CONFIRMA
	LABEL_CONFIRMA2

	LABEL_OFFLINE
)

type IPOctets [4]int

func (display *SerialDisplay) DrawScreen(code string) {

	display.Forth.Send(code + " 0 API") // draw opcode
}

func (display *SerialDisplay) ScreenTags(nome, commVerif int, tags, tagsUnicas ForthNumber) {

	display.DrawScreen(
		fmt.Sprintf(
			"%d lbl %d num"+
				" %d lbl"+
				" %d %d fnm"+ // Tags Val+Mag

				" %d lbl"+
				" %d %d fnm"+ // TagsUnicas Val+Mag

				" %d lbl %d val",

			flick.PORTAL, nome,
			flick.REGIST, tags.Value, tags.Magnitude,
			flick.UNICAS, tagsUnicas.Value, tagsUnicas.Magnitude,
			LABEL_OFFLINE, 6,
		),
	)
}

func (display *SerialDisplay) ScreenAddr(nome, commVerif int, ip IPOctets, leitorOk int) {

	log.Println(ip)

	display.DrawScreen(
		fmt.Sprintf(
			"%d lbl %d num"+
				" %d lbl %d %d %d %d ip"+
				" %d lbl %d val"+
				" %d lbl %d val",

			flick.PORTAL, nome,
			flick.IP, ip[0], ip[1], ip[2], ip[3],
			flick.LEITOR, leitorOk,
			LABEL_OFFLINE, 6,
		),
	)
}

func (display *SerialDisplay) ScreenWifi(nome, commVerif, wifiVerif, LTE4GVerif int, wifiPing int64) {

	display.DrawScreen(
		fmt.Sprintf(
			"%d lbl %d num"+
				" %d lbl %d val"+
				" %d lbl %d val"+
				" %d lbl %d ms",

			flick.PORTAL, nome      ,
			flick.WIFI  , wifiVerif ,
			flick.LTE4G , LTE4GVerif,
			flick.PING  , wifiPing  ,
		),
	)
}

func (display *SerialDisplay) ScreenStat(nome, commVerif int, a1, a2, a3, a4 ForthNumber) {

	display.DrawScreen(
		fmt.Sprintf(
			"%d lbl %d num"+
				" %d %d"+ // A4 Val+Mag
				" %d %d"+ // A3 Val+Mag
				" %d %d"+ // A2 Val+Mag
				" %d %d atn"+ // A1 Val+Mag then display
				" %d lbl %d val",

			flick.PORTAL , nome        ,
			a4.Value     , a4.Magnitude,
			a3.Value     , a3.Magnitude,
			a2.Value     , a2.Magnitude,
			a1.Value     , a1.Magnitude,
			LABEL_OFFLINE, 6           ,
		),
	)
}

func (display *SerialDisplay) ScreenTime(nome, commVerif int) {

	now := time.Now().In(c.ProgramTimezone)
	y, m, d := now.Date()
	//log.Println("now", now)

	display.DrawScreen(
		fmt.Sprintf(
			"%d lbl %d num"+

				// display Time label
				" tim"+

				// Hours, -3 cuz we at GMT-3
				" %d "+

				// Minutes, Seconds
				" %d %d hms"+

				// skip line
				" 22 lbl %d %d %d $DA7E ip"+

				" %d lbl %d val",

			flick.PORTAL , nome                      ,
			now.Hour()   , now.Minute(), now.Second(),
			d            , m           , y           ,
			LABEL_OFFLINE, 6                         ,
		),
	)
}

func (display *SerialDisplay) ScreenUSB(nome, commVerif int, devVerif int) {

	display.DrawScreen(
		fmt.Sprintf(
			"%d lbl %d num"+
				" usb %d val"+
				" fwd"+
				" %d lbl %d val",

			flick.PORTAL ,nome,
			devVerif          ,
			LABEL_OFFLINE, 6  ,
		),
	)
}

func (display *SerialDisplay) ScreenInfoEquip(nome int) {

	display.DrawScreen(
		fmt.Sprintf(
			"%d lbl %d num"+
				// ( ( CA: chafon, FF: impinj ) << 2 ) | ( reader name >> 1 )
				" %d lbl $%s 16 fnm"+
				" %d lbl %d num"+
				" %d lbl $%s hex",

			flick.PORTAL, nome,
			LABEL_RFID  , c.Reader,
			LABEL_SERIE , c.Serie,
			LABEL_SIST  , c.Version,
		),
	)
}

func (display *SerialDisplay) ScreenConfirma() {

	display.DrawScreen(
		fmt.Sprintf(
			"fwd"+
				" %d lbl fwd %d lbl fwd fwd",

			LABEL_CONFIRMA,
			LABEL_CONFIRMA2,
		),
	)
}

func (display *SerialDisplay) ScreenProgress() {

	display.DrawScreen(
		fmt.Sprintf(
			"fwd fwd"+
				" %d lbl fwd fwd",

			LABEL_PROGRESSO,
		),
	)
}

func (display *SerialDisplay) ScreenErr() {

	display.DrawScreen(
		fmt.Sprintf(
			"fwd"+
				" %d lbl fwd %d lbl fwd fwd",

			LABEL_ERRO,
			LABEL_ERRO2,
		),
	)
}
