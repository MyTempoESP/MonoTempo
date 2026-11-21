package usb

import (
	"fmt"
	"path/filepath"
	"strings"
)

func CheckUSBStorageDevice(fs FileSystem) (device string, check bool, err error) {

	const sysBlockPath = "/sys/block/"

	// Read the contents of /sys/block
	blockDevices, err := fs.ReadDir(sysBlockPath)

	if err != nil {

		err = fmt.Errorf("failed to read %s: %w", sysBlockPath, err)

		return
	}

	for _, dev := range blockDevices {

		var realPath string

		devicePath := filepath.Join(sysBlockPath, dev.Name())
		deviceFile := filepath.Join(devicePath, "device")
		realPath, err = fs.EvalSymlinks(deviceFile)

		// USB devices will have "usb" in their symlink path
		if err == nil && strings.Contains(realPath, "/usb") {

			device = "/dev/" + dev.Name()
			check = true

			return
		}
	}

	return
}
