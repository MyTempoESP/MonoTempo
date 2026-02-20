package main

import (
	"fmt"
	"log"
	"os/exec"
)

/*
- Info            : START: Reset visual tag data

- Network         :

- Network Mgmt    : START: Issue a reconnection

- USB Config      : START: Create a report on the USB device.

- System          : START: Fetch and Install the latest version from github.

- AutoUp          : START: Toggle auto-upload mode.

- Upload          : START: Upload all backups.

- #15 (Erase data): START: Erase all data from the device.

- #15 (Shutdown)  : START: Shutdown the device.

#define INFORM_SCREEN 0

#define NETWRK_SCREEN 1

#define NETCFG_SCREEN 2

#define USBCFG_SCREEN 3

#define DATTME_SCREEN 4

#define SYSTEM_SCREEN 5

#define AUTOUP_SCREEN 7

#define UPLOAD_SCREEN 8

#define DELETE_SCREEN 10

#define SHTDWN_SCREEN 11
*/
const (
	infoAction = iota
	antennaAction
	networkAction
	networkMgmtAction
	usbcfgAction
	datetimeAction
	updateAction
	autouploadAction
	uploadAction
	validateAction
	eraseAction
	shutdownAction
)

func CMD(s string) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' > /var/monotempo-data/sig-upload-data", s))
	err := cmd.Run()
	log.Println(err)
}

func AUX(s string) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' > /var/monotempo-data/sig-device-operation", s))
	err := cmd.Run()
	log.Println(err)
}

func PCShutdown()       { CMD("poweroff") }
func PCReboot()         { CMD("reboot") }
func ToggleAutoUpload() { CMD("auto_up") }
func UploadData()       { CMD("normal") }
func SendValidation()   { CMD("validate") }
func PCUpdate()         { CMD("update") }
func CreateUSBReport()  { CMD("stats") }
func FullReset()        { CMD("reset") }
func ResetWifi()        { AUX("reset") }
func Refresh()          { CMD("fatal") }
func Reset4g()          { AUX("lte4g") }
func CopyToUSB()        { CMD("save") }
