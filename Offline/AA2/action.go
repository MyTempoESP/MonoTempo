package main

import (
	"fmt"
	"log"

	"os/exec"
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

func PCReboot()           { CMD("reboot") }
func UploadData()         { CMD("normal") }
func UploadBackup()       { CMD("backup") }
func AtualizarEquip()     { CMD("update") }
func CreateUSBRelatorio() { CMD("stats") }
func ResetarTudo()        { CMD("reset") }
func ResetWifi()          { AUX("reset") }
func Reset4g()            { AUX("lte4g") }
func CopyToUSB()          { CMD("save") }
