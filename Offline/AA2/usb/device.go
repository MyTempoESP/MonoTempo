package usb

import (
	"errors"
	"log"
)

type Device struct {
	Name      string
	FS        FileSystem
	IsMounted bool
}

func (d *Device) Mount(mountPoint string) (err error) {

	log.Println("@ USB - MOUNTING DEVICE")

	if d.IsMounted {

		err = errors.New("Device already mounted!")

		return
	}

	err = Mount(d.Name, mountPoint)

	if err != nil {

		return
	}

	log.Println("@ USB - DEVICE MOUNTED")
	d.IsMounted = true

	return
}

func (d *Device) Umount() (err error) {

	log.Println("@ USB - UMOUNTING DEVICE")

	d.IsMounted = false

	err = Umount()

	return
}

func (d *Device) Check() (check bool, err error) {

	d.Name, check, err = CheckUSBStorageDevice(d.FS)

	if err != nil {

		return
	}

	if !check && d.IsMounted {

		log.Println("@ USB - DEVICE REMOVED")

		d.IsMounted = false
		err = Umount()

		return
	}

	return
}
