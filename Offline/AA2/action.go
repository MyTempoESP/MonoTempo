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

- Upload          : START: Upload all tag data currently stored + pending.

- Upload (Backup) : START: Upload all backups.

- #15 (Erase data): START: Erase all data from the device.

- #15 (Shutdown)  : START: Shutdown the device.

#define INFORM_SCREEN 0

#define NETWRK_SCREEN 1

#define NETCFG_SCREEN 2

#define USBCFG_SCREEN 3

#define DATTME_SCREEN 4

#define SYSTEM_SCREEN 5

#define UPLOAD_SCREEN 6

#define BACKUP_SCREEN 7

#define DELETE_SCREEN 8

#define SHTDWN_SCREEN 9
*/
const (
	INFO_ACTION = iota
	NETWORK_ACTION
	NETWORK_MGMT_ACTION
	USBCFG_ACTION
	DATETIME_ACTION
	UPDATE_ACTION
	UPLOAD_ACTION
	UPLOAD_BACKUP_ACTION
	ERASE_ACTION
	SHUTDOWN_ACTION
)

// CMD hangs the screen until a restart issued by chancelor itself
func CMD(s string) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' > /var/monotempo-data/sig-upload-data", s))
	err := cmd.Run()
	log.Println(err)

	select {}
}

func AUX(s string) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' > /var/monotempo-data/sig-device-operation", s))
	err := cmd.Run()
	log.Println(err)
}

// Shuts down the PC. FIXME: Implement poweroff command in chancelor.
func PCShutdown() { CMD("poweroff") }

// Reboots the PC.
func PCReboot() { CMD("reboot") }

// Uploads normal tag data.
func UploadData() { CMD("normal") }

// Uploads backup tag data.
func UploadBackup() { CMD("backup") }

// Updates the equipment to the latest version.
func PCUpdate() { CMD("update") }

// Creates a USB report with device statistics.
func CreateUSBReport() { CMD("stats") }

// Resets all device data.
func FullReset() { CMD("reset") }

// Resets the Wi-Fi configuration.
func ResetWifi() { AUX("reset") }

// Refreshes the system state, used for fatal errors.
func Refresh() { CMD("fatal") }

// Resets the 4G LTE configuration.
func Reset4g() { AUX("lte4g") }

// Copies data to a USB device.
func CopyToUSB() { CMD("save") }
