package main

import (
	"log"

	"os/exec"
)

func ResetWifi() {
	cmd := exec.Command("sh", "-c", "nmcli conn down MyTempo; nmcli conn up MyTempo")
	go func() { log.Println(cmd.Run()) }()
}

func ResetarTudo() {
	cmd := exec.Command("sh", "-c", "echo 'reset' > /var/monotempo-data/sig-upload-data")
	err := cmd.Run()
	log.Println(err)
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

func CopyToUSB() {
	cmd := exec.Command("sh", "-c", "echo 'save' > /var/monotempo-data/sig-upload-data")
	err := cmd.Run()
	log.Println(err)
}
