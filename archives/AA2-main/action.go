package main

import (
	"fmt"
	"time"

	file "aa2/file"
	usb "aa2/usb"
)

func CopyToUSB(device *usb.Device, file *file.File) (err error) {

	if !device.IsMounted {

		err = device.Mount("/mnt")

		if err != nil {

			return
		}
	}

	now := time.Now()

	err = file.Upload(fmt.Sprintf("/mnt/MYTEMPO-%02d_%02d_%02d", now.Hour(), now.Minute(), now.Second()))

	return
}
