package file

import (
	"os/exec"
)

func copyFile(src, dst string) (err error) {

	args := []string{src, dst}
	cmd := exec.Command("cp", args...)

	err = cmd.Run()

	return
}
