package main

import (
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"aa2/flick"
	"aa2/intSet"
	"aa2/lcdlogger"
	"aa2/pinger"
	"aa2/usb"
)

func PopulateTagSet(tagSet *intSet.IntSet) {

	b, err := os.ReadFile("/var/monotempo-data/TAGS")

	if err != nil {

		return
	}

	for _, s := range strings.Split(string(b), "\n") {

		n, err := strconv.Atoi(s)

		if err != nil {

			continue
		}

		tagSet.Insert(n)
	}
}

func (a *Ay) Process() {

	var (
		tags     atomic.Int64
		tagsUSB  atomic.Int64
		antennas [4]atomic.Int64
		tagSet   intSet.IntSet
	)

	tagSet = intSet.New()

	PopulateTagSet(&tagSet)

	tags_start_at := os.Getenv("TAG_COUNT_START_AT")

	go func() {

		if tags_start_at != "" {

			tags_start_at, err := strconv.Atoi(tags_start_at)

			if err == nil {

				tags.Store(int64(tags_start_at))
			}
		}

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
		}
	}()

	var device = usb.Device{}
	device.Name = "/dev/sdb"
	device.FS = usb.OSFileSystem{}

	var readerIP = os.Getenv("READER_IP")
	// var readerOctets = lcdlogger.IPIfy(readerIP)
	var readerState atomic.Bool

	var netPing atomic.Int64
	var netState atomic.Bool
	//var readerPing atomic.Int64

	ReaderPinger := pinger.NewPinger(readerIP, &readerState, nil)
	WifiPinger := pinger.NewPinger("mytempo.esp.br", &netState, &netPing)

	go ReaderPinger.Run()
	go WifiPinger.Run()

	display, displayErr := lcdlogger.NewSerialDisplay()

	if displayErr != nil {

		return
	}

	NUM_EQUIP, err := strconv.Atoi(os.Getenv("MYTEMPO_DEVID"))

	if err != nil {
		return
	}

	go func() {

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

				wiok := flick.OK
				if !netState.Load() {

					wiok = flick.X
				}

				display.ScreenWifi(
					NUM_EQUIP,
					commVerif,
					/* leitor OK? */ ok,
					/* internet OK? */ wiok,
					netPing.Load(),
				)
			case lcdlogger.SCREEN_4G:

				ok := flick.OK

				if !readerState.Load() {

					ok = flick.DESLIGAD
				}

				wiok := flick.OK
				if !netState.Load() {

					wiok = flick.X
				}

				display.Screen4g(
					NUM_EQUIP,
					commVerif,
					/* leitor OK? */ ok,
					/* internet OK? */ wiok,
					netPing.Load(),
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
			case lcdlogger.SCREEN_UPLOAD:
				display.ScreenConfirmaUpload()
			case lcdlogger.SCREEN_UPLOAD_BACKUP:
				display.ScreenConfirmaUploadBackup()
			}

			display.SwitchScreens()

			if action, hasAction := display.Action(); hasAction {

				switch action {
				case lcdlogger.ACTION_RESET:
					display.ScreenConfirma()
				case lcdlogger.ACTION_UPLOAD:
					fallthrough
				case lcdlogger.ACTION_UPLOAD_BACKUP:
					display.ScreenUpload()
				default:
					display.ScreenProgress()
				}

				err = nil

				switch action {
				case lcdlogger.ACTION_UPLOAD:
					UploadData()
					select {}
				case lcdlogger.ACTION_UPLOAD_BACKUP:
					UploadBackup()
					select {}
					/*				case lcdlogger.ACTION_REBOOT:
									PCReboot()
									select {}
					*/
				case lcdlogger.ACTION_USB:
					CopyToUSB()
					select {}
				case lcdlogger.ACTION_WIFI_RESET:
					{
						ResetWifi()

						WifiPinger.Stop()
						WifiPinger = pinger.NewPinger("mytempo.esp.br", &netState, &netPing)

						go WifiPinger.Run()

						<-time.After(2 * time.Second)
					}
				case lcdlogger.ACTION_RESET:
					{
						hasKey := display.WaitKeyPress(5 * time.Second)

						if !hasKey { // timeout

							continue
						}

						display.ScreenProgress()

						ResetarTudo()
						select {}
					}
				case lcdlogger.ACTION_TAGS:
					{

						for i := range 4 {
							antennas[i].Store(0)
						}

						tags.Store(0)
						tagSet.Clear()
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
