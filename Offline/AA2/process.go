package main

import (
	"log"
	"os"
	"sync/atomic"
	"time"

	"aa2/file"
	"aa2/flick"
	"aa2/intSet"
	"aa2/lcdlogger"
	"aa2/pinger"
	"aa2/usb"
)

func (a *Ay) Process() {

	var (
		tags     atomic.Int64
		tagsUSB  atomic.Int64
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
			tagsUSB.Add(1)

			tagSet.Insert(t.Epc)

			//001000000000000036513:34:21.097
			tagsFile.Insert(t.FormatoRefinado)
		}
	}()

	var device = usb.Device{}

	device.Name = "/dev/sdb"
	device.FS = usb.OSFileSystem{}

	var readerIP = os.Getenv("READER_IP")
	var readerOctets = lcdlogger.IPIfy(readerIP)
	var readerState atomic.Bool
	//var readerPing atomic.Int64

	go pinger.NewPinger(readerIP, &readerState, nil)

	display, displayErr := lcdlogger.NewSerialDisplay()

	if displayErr != nil {

		return
	}

	go func() {

		const NUM_EQUIP = 501

		for {

			commVerif := flick.DESLIGAD

			switch display.Screen {
			case lcdlogger.SCREEN_TAGS:
				display.ScreenTags(
					NUM_EQUIP,
					commVerif,
					/* Tags    */ lcdlogger.ToForthNumber(tags.Load()),
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
					/* IP         */ readerOctets,
					/* leitor OK? */ ok,
				)
			case lcdlogger.SCREEN_WIFI:

				wifi := flick.DESLIGAD
				lte := flick.DESLIGAD

				display.ScreenWifi(
					NUM_EQUIP,
					commVerif,
					/* WIFI */ wifi,
					/* 4G   */ lte,
					-1,
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
			case lcdlogger.SCREEN_INFO_EQUIP:
				display.ScreenInfoEquip(NUM_EQUIP)
			}

			display.SwitchScreens()

			if action, hasAction := display.Action(); hasAction {

				switch action {
				case lcdlogger.ACTION_RESET:
					display.ScreenConfirma()
				default:
					display.ScreenProgress()
				}

				err = nil

				switch action {
				case lcdlogger.ACTION_WIFI:
					UploadData()
					select {}
				case lcdlogger.ACTION_TIME: /* empty */
					UploadBackup()
					select {}
				case lcdlogger.ACTION_RESET:
					{
						hasKey := display.WaitKeyPress(5 * time.Second)

						if !hasKey { // timeout

							continue
						}

						display.ScreenProgress()

						// resetar equip

						// err = ResetarTudo()
						tagsFile.Clear()
					}
					fallthrough // resetar tags
				case lcdlogger.ACTION_TAGS:
					{

						for i := range 4 {
							antennas[i].Store(0)
						}

						tags.Store(0)
						tagSet.Clear()
					}
				case lcdlogger.ACTION_USB:
					{

						err = CopyToUSB(&device, &tagsFile)

						if err == nil {

							tagsUSB.Load()

							<-time.After(5 * time.Second)

							tagsUSB.Store(0)
						}
					}
				default:
					continue // no action
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
