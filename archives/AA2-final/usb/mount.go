package usb

import (
	"os/exec"
)

func Mount(device, mountPoint string) (err error) {

	args := []string{device + "1", mountPoint}
	cmd := exec.Command("mount", args...)
	err = cmd.Run()

	return
}

func Umount() (err error) {

	cmd := exec.Command("umount", "/mnt")
	err = cmd.Run()

	return
}

