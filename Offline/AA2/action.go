package main

import (
	"fmt"
	"log"
	"time"

	c "aa2/constant"
	file "aa2/file"
	usb "aa2/usb"
	"os/exec"
)

func ResetarTudo() (err error) {

	// delete entire times database

	return
}

func PCReboot() {
	cmd := exec.Command("sh", "-c", "echo 'reboot' > /var/monotempo-data/sig-upload-data")
	err := cmd.Run()
	log.Println(err)
}

func UploadData() {
	cmd := exec.Command("sh", "-c", "echo 'normal' > /var/monotempo-data/sig-upload-data")
	err := cmd.Run()
	log.Println(err)
}

func UploadBackup() {
	cmd := exec.Command("sh", "-c", "echo 'backup' > /var/monotempo-data/sig-upload-data")
	err := cmd.Run()
	log.Println(err)
}

func CopyToUSB(device *usb.Device, file *file.File) (err error) {

	if !device.IsMounted {

		err = device.Mount("/mnt")

		if err != nil {

			log.Println("Error mounting")

			return
		}
	}

	now := time.Now().In(c.ProgramTimezone)

	log.Println("copying")

	err = file.Upload(fmt.Sprintf("/mnt/MYTEMPO-%02d_%02d_%02d", now.Hour(), now.Minute(), now.Second()))

	if err != nil {

		log.Println(err)

		return
	}

	err = device.Umount()

	return
}
