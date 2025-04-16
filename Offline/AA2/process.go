package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"aa2/com"
	"aa2/constant"
	"aa2/intSet"
	"aa2/pinger"
	"aa2/usb"
)

func countDir(path string) (n int, err error) {

	f, err := os.Open(path)

	if err != nil {

		return
	}

	list, err := f.Readdirnames(-1)

	f.Close()

	if err != nil {

		return
	}

	n = len(list)

	return
}

func populateTagSet(tagSet *intSet.IntSet, permanentSet *intSet.IntSet) {

	b, err := os.ReadFile("/var/monotempo-data/TAGS")

	if err != nil {

		return
	}

	for s := range strings.SplitSeq(string(b), "\n") {

		n, err := strconv.Atoi(s)

		if err != nil {

			continue
		}

		tagSet.Insert(n)
		permanentSet.Insert(n)
	}
}

func checkAction(actionString string, tagSet *intSet.IntSet, tags *atomic.Int64, antennas *[4]atomic.Int64) {

	idx := strings.Index(actionString, "$")

	if idx == -1 {
		return
	}

	actionString = actionString[idx:]

	if strings.HasPrefix(actionString, "$MYTMP;") {
		actionString = strings.TrimPrefix(actionString, "$MYTMP;")
		action, err := strconv.Atoi(strings.TrimSpace(actionString))

		if err != nil {
			return
		}

		switch action {
		case INFO_ACTION:
			tagSet.Clear()
			tags.Store(0)
			antennas[0].Store(0)
			antennas[1].Store(0)
			antennas[2].Store(0)
			antennas[3].Store(0)
		case NETWORK_ACTION:
		case NETWORK_MGMT_ACTION:
			ResetWifi()
		case DATETIME_ACTION:

		// these actions hang
		case USBCFG_ACTION:
			CreateUSBReport()
			select {}
		case UPDATE_ACTION:
			PCUpdate()
			select {}
		case UPLOAD_ACTION:
			UploadData()
			select {}
		case UPLOAD_BACKUP_ACTION:
			UploadBackup()
			select {}
		case ERASE_ACTION:
			FullReset()
			select {}
		case SHUTDOWN_ACTION:
			PCShutdown()
			select {}
		default:
			return
		}

		<-time.After(time.Second * 3)
	}
}

const (
	STATE_TAG_REPORT = iota
	STATE_ANTENNA_REPORT
	STATE_PC_DATA_REPORT
)

// function for the state transition, it goes: 0, 0, 0, 1, 2, 0 ...
func transitionStep(c int) int {
	r := c % 5

	if r < 3 {
		return 0
	} else {
		return r - 2
	}
}

func (a *Ay) Process() {

	var (
		pcData          *com.PCData = &com.PCData{}
		tagsUSB         atomic.Int64
		tagSet          intSet.IntSet = intSet.New()
		permanentTagSet intSet.IntSet = intSet.New()
	)

	populateTagSet(&tagSet, &permanentTagSet)

	tags_start_at := os.Getenv("TAG_COUNT_START_AT")

	go func() {

		if tags_start_at != "" {

			tags_start_at, err := strconv.Atoi(tags_start_at)

			if err == nil {

				pcData.Tags.Store(int64(tags_start_at))
			}
		}

		for t := range a.Tags {

			if t.Antena == 0 {

				/*
					Antena 0 não exist
				*/

				continue
			}

			pcData.Antennas[(t.Antena-1)%4].Add(1)

			pcData.Tags.Add(1)
			tagsUSB.Add(1)

			tagSet.Insert(t.Epc)
			permanentTagSet.Insert(t.Epc)
		}
	}()

	// Inicializa o SerialSender com uma taxa de baud de 115200
	sender, err := com.NewSerialSender(115200, constant.SerialPortOverride)

	if err != nil {
		log.Printf("Falha ao inicializar o SerialSender: %v", err)
		return
	}

	// I AM DUMB AS FUCK
	// defer sender.Close()

	var device = usb.Device{}
	device.Name = "/dev/sdb"
	device.FS = usb.OSFileSystem{}

	var readerIP = os.Getenv("READER_IP")

	go pinger.NewJSONPinger(&pcData.CommStatus)

	ReaderPinger := pinger.NewPinger(readerIP, &pcData.RfidStatus, nil)

	go ReaderPinger.Run()

	sysver, err := strconv.Atoi(constant.VersionNum)

	if err != nil {
		sysver = 0
		log.Printf("Falha ao converter a versão do sistema: %v, utilizando 0", err)
	}

	pcData.Tags.Store(0)
	pcData.UniqueTags.Store(0)

	pcData.SysVersion = sysver

	backupDirs, err := countDir("/var/monotempo-data/backup")

	if err != nil {
		log.Printf("Erro ao listar diretórios de backup: %v", err)
		pcData.Backups = 0
	} else {
		pcData.Backups = backupDirs
	}

	deviceId, err := strconv.Atoi(constant.DeviceId)
	if err != nil {
		log.Printf("Erro ao converter o hostname para número: %v", err)
		pcData.SysCodeName = 500
	} else {
		pcData.SysCodeName = deviceId
	}

	// Envia os dados iniciais
	pcData.SendPCDataReport(sender)
	<-time.After(time.Second * 3)

	//NUM_EQUIP, err := strconv.Atoi(os.Getenv("MYTEMPO_DEVID"))

	// TODO: revert everything you did today again

	go func() {

		switcherTicker := time.NewTicker(1 * time.Second)
		sendTicker := time.NewTicker(120 * time.Millisecond)
		state := STATE_TAG_REPORT

		// step counter for state transitions
		c := 0

		for range sendTicker.C {

			pcData.UniqueTags.Store(int32(tagSet.Count()))

			pcData.PermanentUniqueTags.Store(int32(permanentTagSet.Count()))

			usbOk, _ := device.Check()
			pcData.UsbStatus.Store(usbOk)

			switch state {
			case STATE_TAG_REPORT:
				pcData.SendTagReport(sender)
			case STATE_ANTENNA_REPORT:
				if constant.ReaderType == "impinj" {
					pcData.SendAntennaReport(sender)
				} else {
					pcData.SendPCDataReport(sender)
				}
			case STATE_PC_DATA_REPORT:
				pcData.SendPCDataReport(sender)
			}

			actionString, hasAction := sender.Recv()

			if hasAction {
				checkAction(actionString, &tagSet, &pcData.Tags, &pcData.Antennas)
			}

			select {
			case <-switcherTicker.C:
				c += 1 // step
				state = transitionStep(c)
			default:
			}
		}
	}()
}
