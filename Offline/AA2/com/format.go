package com

import (
	"aa2/logparse"
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

type PCData struct {
	Tags                atomic.Int64
	UniqueTags          atomic.Int32
	CommStatus          atomic.Bool
	WifiStatus          atomic.Bool
	Lte4Status          atomic.Bool
	RfidStatus          atomic.Bool
	UsbStatus           atomic.Bool
	AutoUploadStatus    atomic.Bool
	PermanentUniqueTags atomic.Int32
	Antennas            [4]AntennaSignal

	// constants
	SysVersion  int
	SysCodeName int
	Backups     int
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

/*
This is based on the following checksum function

	bool check_sum(SafeString &msg) {
	  int idxStar = msg.indexOf('*');

	  cSF(check_sum_hex, 2);

	  msg.substring(check_sum_hex, idxStar + 1);

	  long sum = 0;

	  if (!check_sum_hex.hexToLong(sum)) {
	    return false;
	  }

	  for (size_t i = 1; i < idxStar; i++) {
	    sum ^= msg[i];
	  }

	  return (sum == 0);
	}
*/
func withChecksum(data string) string {
	var checksum byte

	for i := range len(data) {
		checksum ^= data[i]
	}

	return fmt.Sprintf("$%s*%02X", data, checksum)
}

func epoch() int64 {
	unix := time.Now().Unix()
	_, offset := time.Now().Zone()
	return unix + int64(offset)
}

func (pd *PCData) formatPCDataReport() string {
	currentEpoch := epoch()

	f := fmt.Sprintf("MYTMP;%d;%d;P;%d;%d;%d;%d;%d;%d;%d;%d;%d",
		pd.Tags.Load(), pd.UniqueTags.Load(), boolToInt(pd.CommStatus.Load()),
		boolToInt(pd.RfidStatus.Load()), boolToInt(pd.UsbStatus.Load()),
		boolToInt(pd.AutoUploadStatus.Load()),
		pd.SysVersion, pd.SysCodeName, pd.Backups, pd.PermanentUniqueTags.Load(),
		currentEpoch)

	return withChecksum(f)
}

func (pd *PCData) formatAntennaReport() string {
	currentEpoch := epoch()

	f := fmt.Sprintf("MYTMP;%d;%d;A;%d;%d;%d;%d;%d",
		pd.Tags.Load(), pd.UniqueTags.Load(), pd.Antennas[0].Level(),
		pd.Antennas[1].Level(), pd.Antennas[2].Level(), pd.Antennas[3].Level(), currentEpoch)

	return withChecksum(f)
}

func (pd *PCData) formatTagReport() string {
	currentEpoch := epoch()

	f := fmt.Sprintf("MYTMP;%d;%d;T;%d",
		pd.Tags.Load(), pd.UniqueTags.Load(), currentEpoch)

	return withChecksum(f)
}

func (pd *PCData) formatLogReport(equipStatus *logparse.EquipStatus) string {
	currentEpoch := epoch()

	f := fmt.Sprintf("MYTMP;0;0;L;%d;%d;%d;%d;%d;%d",
		equipStatus.UploadCount, equipStatus.Databases, equipStatus.AvgProctime,
		equipStatus.Errcount, boolToInt(equipStatus.Status), currentEpoch)

	return withChecksum(f)
}

func (pd *PCData) SendTagReport(sender *SerialSender) {
	data := pd.formatTagReport()
	log.Println("Sending TagReport:", data)
	sender.SendData(data)
}

func (pd *PCData) SendAntennaReport(sender *SerialSender) {
	data := pd.formatAntennaReport()
	log.Println("Sending AntennaReport:", data)
	sender.SendData(data)
}

func (pd *PCData) SendPCDataReport(sender *SerialSender) {
	data := pd.formatPCDataReport()
	log.Println("Sending PCDataReport:", data)
	sender.SendData(data)
}

func (pd *PCData) SendLogReport(sender *SerialSender, equipStatus *logparse.EquipStatus) {
	data := pd.formatLogReport(equipStatus)
	log.Println("Sending LogReport:", data)
	sender.SendData(data)
}
