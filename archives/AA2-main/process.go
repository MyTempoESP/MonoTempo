package main

import (
	"log"
	"os"
	"sync/atomic"
	"time"

	"aa2/file"
	"aa2/intSet"
	"aa2/lcdlogger"
	"aa2/pinger"
	"aa2/usb"
	"github.com/MyTempoesp/flick"
)

func (a *Ay) Process() {

	var (
		tags     atomic.Int64
		antennas [4]atomic.Int64
	)

	tagSet := intSet.New()

	tagsFile, err := file.NewFile("tags")

	if err != nil {

		log.Println(err)
	}

	go func() {

		for t := range a.Tags {

			if t.Antena == 0 {
				/*
					Antena 0 não exist
				*/

				continue
			}

			antennas[(t.Antena-1)%4].Add(1)

			tags.Add(1)

			tagSet.Insert(t.Epc)

			//001000000000000036513:34:21.097
			tagsFile.Insert(t.FormatoRefinado)
		}
	}()

	/*
		Resetting the tags:

		TODO: add a clear method of
		resetting tags on the database

		tagSet.Clear()
		tags.Store(0)
		for i := range 4 {

			antennas[i].Store(0)
		}

	*/

	var device = usb.Device{}

	device.Name = "/dev/sdb"
	device.FS = usb.OSFileSystem{}

	var readerIP = os.Getenv("READER_IP")
	var readerOctets = lcdlogger.IPIfy(readerIP)
	var readerState atomic.Bool
	//var readerPing atomic.Int64

	go pinger.NewPinger(readerIP, &readerState, nil)

	var wifiIP = os.Getenv("MYTEMPO_API_URL")
	var wifiState atomic.Bool
	var wifiPing atomic.Int64

	go pinger.NewPinger(wifiIP, &wifiState, &wifiPing)

	display, displayErr := lcdlogger.NewSerialDisplay()

	/* > Monitoring can be skipped if NewSerialDisplay() errors out, disabling the routine in Line 221 */
	if displayErr != nil {

		return
	}

	go func() {

		const NUM_EQUIP = 701

		for {

			commVerif := flick.WEB

			switch display.Screen {
			case lcdlogger.SCREEN_TAGS:
				display.ScreenTags(
					NUM_EQUIP,
					commVerif,
					/* Tags */ lcdlogger.ToForthNumber(tags.Load()),
					/* Atletas */ lcdlogger.ToForthNumber(tagSet.Count()),
				)
			case lcdlogger.SCREEN_ADDR:

				ok := flick.OK

				if !readerState.Load() {

					ok = flick.DESLIGAD
				}

				display.ScreenAddr(
					NUM_EQUIP,
					commVerif,
					/* IP */ readerOctets,
					/* leitor OK? */ ok,
				)
			case lcdlogger.SCREEN_WIFI:
				display.ScreenWifi(
					NUM_EQUIP,
					commVerif,
					/* WIFI */ flick.CONECTAD,
					/* 4G */ flick.DESLIGAD,
					wifiPing.Load(),
				)
			case lcdlogger.SCREEN_STAT:
				display.ScreenStat(
					NUM_EQUIP,
					commVerif,
					lcdlogger.ToForthNumber(antennas[0].Load()),
					lcdlogger.ToForthNumber(antennas[1].Load()),
					lcdlogger.ToForthNumber(antennas[2].Load()),
					lcdlogger.ToForthNumber(antennas[3].Load()),
				)
			case lcdlogger.SCREEN_TIME:
				display.ScreenTime(
					NUM_EQUIP,
					commVerif,
				)
			case lcdlogger.SCREEN_USB:
				devCheck, err := device.Check()

				if err != nil {

					continue
				}

				devVerif := flick.X

				if devCheck {

					devVerif = flick.CONECTAD
				}

				display.ScreenUSB(
					NUM_EQUIP,
					commVerif,
					devVerif,
				)
			}

			display.HandleActionButton()
			display.SwitchScreens()

			if action, hasAction := display.Action(); hasAction {

				display.ScreenProgress()

				err = nil

				switch action {
				case lcdlogger.ACTION_WIFI: /* empty */
				case lcdlogger.ACTION_TIME: /* empty */
				case lcdlogger.ACTION_TAGS:
					tags.Store(0)
					tagSet.Clear()
				case lcdlogger.ACTION_USB:
					err = CopyToUSB(&device, &tagsFile)
				}

				<-time.After(1 * time.Second) // min 1 sec

				if err != nil {

					display.ScreenErr()

					<-time.After(5 * time.Second)

					continue
				}
			}

			display.HandleAltActionButton()
			display.ProcessAltAction()

			if action, hasAction := display.AltAction(); hasAction {

				display.ScreenProgress()

				err = nil

				switch action {
				case lcdlogger.ALT_ACTION_RESTART:
					err = RestartComputer()
				}

				<-time.After(1 * time.Second) // min 1 sec

				if err != nil {

					display.ScreenErr()

					<-time.After(5 * time.Second)

					continue
				}
			}

			time.Sleep(50 * time.Millisecond)
		}
	}()
}
