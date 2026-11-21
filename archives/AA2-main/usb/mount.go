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

func Umount() {

	cmd := exec.Command("umount", "/mnt")
	cmd.Run()
}

