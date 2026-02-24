package main

import (
	"aa2/com"
	"aa2/constant"
	"aa2/intSet"
	"aa2/logparse"
	"aa2/pinger"
	"aa2/usb"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
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

func checkAction(actionString string, state *int, smoother *TagSmoother, tagSet *intSet.IntSet, permanentTagSet *intSet.IntSet, tags *atomic.Int64, antennas *[4]com.AntennaSignal) {
	idx := strings.Index(actionString, "$")

	if idx == -1 {
		return
	}

	actionString = actionString[idx:]

	actionString, _ = strings.CutPrefix(actionString, "$MYTMP;")

	action, err := strconv.Atoi(strings.TrimSpace(actionString))
	if err != nil {
		return
	}

	switch action {
	case infoAction:
		if smoother != nil {
			smoother.Clear()
		} else {
			tagSet.Clear()
			permanentTagSet.Clear()
			tags.Store(0)
		}
		for i := range antennas {
			antennas[i].Clear()
		}
	case antennaAction:
		*state = stateAntennaReport
	case networkAction:
	case networkMgmtAction:
		ResetWifi()
		<-time.After(time.Second * 2)
	case datetimeAction:

	// these actions hang
	case usbcfgAction:
		CreateUSBReport()
		select {}
	case updateAction:
		PCUpdate()
		select {}
	case autouploadAction:
		ToggleAutoUpload()
		select {}
	case uploadAction:
		UploadData()
		select {}
	case validateAction:
		SendValidation()
		select {}
	case eraseAction:
		FullReset()
		select {}
	case shutdownAction:
		PCShutdown()
		select {}
	}
}

const (
	stateTagReport = iota
	stateAntennaReport
	statePCDataReport
)

var (
	states   = [...]int{0, 0, 0, 1, 1, 2, 2}
	maxState = len(states)
)

func transitionStep(c int) int {
	return states[c%maxState]
}

func (a *Ay) Process() {
	var (
		pcData  = &com.PCData{}
		tagsUSB atomic.Int64

		tagSet          = intSet.New()
		permanentTagSet = intSet.New()
	)

	populateTagSet(&tagSet, &permanentTagSet)

	tagsStartAt := os.Getenv("TAG_COUNT_START_AT")

	deviceID, err := strconv.Atoi(constant.DeviceID)

	if err != nil {
		a.logger.Error("Erro ao converter o hostname para número", zap.Error(err))
		pcData.SysCodeName = 500
	} else {
		pcData.SysCodeName = deviceID
	}

	var smoother *TagSmoother
	if deviceID >= 700 {
		smoother = newTagSmoother(&pcData.Tags, &pcData.UniqueTags, &pcData.PermanentUniqueTags, &tagSet, &permanentTagSet)
	}

	go func() {
		if tagsStartAt != "" {
			tagsStartAt, err := strconv.Atoi(tagsStartAt)
			if err == nil {
				pcData.Tags.Store(int64(tagsStartAt))
			}
		}

		for t := range a.Tags {

			if t.Antena == 0 {
				continue
			}

			pcData.Antennas[(t.Antena-1)%4].Record()
			tagsUSB.Add(1)

			if smoother != nil {
				smoother.Push(t.Epc)
			} else {
				pcData.Tags.Add(1)
				tagSet.Insert(t.Epc)
				permanentTagSet.Insert(t.Epc)
			}
		}
	}()

	// Inicializa o SerialSender com uma taxa de baud de 115200
	sender, err := com.NewSerialSender(115200, constant.SerialPortOverride, a.logger)
	if err != nil {
		a.logger.Error("Falha ao inicializar o SerialSender",
			zap.Error(err),
		)
		return
	}

	// sender is intentionally never closed; it must outlive the process.

	device := usb.Device{}
	device.Name = "/dev/sdb"
	device.FS = usb.OSFileSystem{}

	readerIP := os.Getenv("READER_IP")

	go pinger.NewJSONPinger(&pcData.CommStatus, a.logger)

	ReaderPinger := pinger.NewPinger(readerIP, &pcData.RfidStatus, nil)

	go ReaderPinger.Run()

	sysver, err := strconv.Atoi(constant.VersionNum)
	if err != nil {
		sysver = 0
		a.logger.Error("Falha ao converter a versão do sistema, utilizando 0", zap.Error(err))
	}

	autoupload, err := strconv.Atoi(constant.AutoUploadEnabled)
	if err != nil {
		autoupload = 0
		a.logger.Error("Falha ao obter o estado do envio automatico, utilizando 0", zap.Error(err))
	}

	pcData.AutoUploadStatus.Store(autoupload != 0)

	pcData.Tags.Store(0)
	pcData.UniqueTags.Store(0)

	pcData.SysVersion = sysver

	backupDirs, err := countDir("/var/monotempo-data/backup")

	if err != nil {
		a.logger.Error("Erro ao listar diretórios de backup", zap.Error(err))
		pcData.Backups = 0
	} else {
		pcData.Backups = backupDirs
	}

	pcData.SendPCDataReport(sender)
	<-time.After(time.Second * 3)

	equipStatus, err := logparse.ParseJSONLog("/var/monotempo-data/logs/pc.log")

	if err == nil {
		for range 5 {
			pcData.SendLogReport(sender, &equipStatus)
			<-time.After(time.Millisecond * 120)
		}
	}

	// NUM_EQUIP, err := strconv.Atoi(os.Getenv("MYTEMPO_DEVID"))
	// TODO: revert everything you did today again

	go func() {
		switcherTicker := time.NewTicker(400 * time.Millisecond)
		sendTicker := time.NewTicker(120 * time.Millisecond)
		state := stateTagReport

		// step counter for state transitions
		c := 0

		for range sendTicker.C {

			if smoother == nil {
				pcData.UniqueTags.Store(int32(tagSet.Count()))
				pcData.PermanentUniqueTags.Store(int32(permanentTagSet.Count()))
			}

			// usbOk, _ := device.Check()
			// pcData.UsbStatus.Store(usbOk)

			switch state {
			case stateTagReport:
				pcData.SendTagReport(sender)
			case stateAntennaReport:
				if constant.ReaderType == "impinj" {
					pcData.SendAntennaReport(sender)
				} else {
					pcData.SendPCDataReport(sender)
				}
			case statePCDataReport:
				pcData.SendPCDataReport(sender)
			}

			actionString, hasAction := sender.Recv()

			if hasAction {
				checkAction(actionString, &state, smoother, &tagSet, &permanentTagSet, &pcData.Tags, &pcData.Antennas)
			}

			select {
			case <-switcherTicker.C:
				c++
				state = transitionStep(c)
			default:
			}
		}
	}()
}
