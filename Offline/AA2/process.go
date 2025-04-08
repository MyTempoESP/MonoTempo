package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"aa2/constant"
	"aa2/intSet"
	"aa2/pinger"
	"aa2/usb"

	com "github.com/TsukiGva2/comunica_serial"
	probing "github.com/prometheus-community/pro-bing"
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

	// Inicializa o SerialSender com uma taxa de baud de 115200
	sender, err := com.NewSerialSender(115200)

	if err != nil {
		log.Printf("Falha ao inicializar o SerialSender: %v", err)
		return
	}

	defer sender.Close()

	var device = usb.Device{}
	device.Name = "/dev/sdb"
	device.FS = usb.OSFileSystem{}

	var readerIP = os.Getenv("READER_IP")
	var readerState atomic.Bool
	var netState atomic.Bool
	var lte4gState atomic.Bool
	var Lte4gPinger *probing.Pinger

	go pinger.NewJSONPinger(&netState)
	ReaderPinger := pinger.NewPinger(readerIP, &readerState, nil)

	go ReaderPinger.Run()
	go func() {
		for {
			Lte4gPinger = pinger.NewPinger("192.168.100.1", &lte4gState, nil)
			Lte4gPinger.Run()
			<-time.After(1 * time.Second)
			log.Println("4gPING STOPPED")
		}
	}()

	sysver, err := strconv.Atoi(constant.VersionNum)

	if err != nil {
		sysver = 0
		log.Printf("Falha ao converter a versão do sistema: %v, utilizando 0", err)
	}

	pcData := &com.PCData{}
	pcData.Tags.Store(0)
	pcData.UniqueTags.Store(0)
	pcData.CommStatus.Store(false)
	pcData.WifiStatus.Store(false)
	pcData.Lte4Status.Store(false)
	pcData.RfidStatus.Store(false)
	pcData.SysVersion.Store(int32(sysver))

	backupDirs, err := os.ReadDir("/var/monotempo/backup")

	if err != nil {
		log.Printf("Erro ao listar diretórios de backup: %v", err)
		pcData.Backups.Store(0)
	} else {
		pcData.Backups.Store(int32(len(backupDirs)))
	}

	// Envia os dados iniciais
	pcData.Send(sender)
	<-time.After(time.Second * 2)

	//NUM_EQUIP, err := strconv.Atoi(os.Getenv("MYTEMPO_DEVID"))

	go func() {

		// Configura um ticker para enviar dados periodicamente
		ticker := time.NewTicker(120 * time.Millisecond)

		defer ticker.Stop()

		for range ticker.C {

			usbOk, _ := device.Check()

			pcData.Tags.Store(tags.Load())
			pcData.UniqueTags.Store(int32(tagSet.Count()))

			pcData.RfidStatus.Store(readerState.Load())
			pcData.Lte4Status.Store(lte4gState.Load())
			pcData.WifiStatus.Store(netState.Load())
			pcData.CommStatus.Store(netState.Load())
			pcData.UsbStatus.Store(usbOk)

			pcData.Send(sender)
		}
	}()
}
