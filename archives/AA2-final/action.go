package main

import (
	"fmt"
	"log"
	"time"

	c "aa2/constant"
	file "aa2/file"
	usb "aa2/usb"
	mysql "github.com/mytempoesp/mysql-easy"
)

func ResetarTudo() (err error) {

	db, err := mysql.ConfiguraDB()

	if err != nil {

		return
	}

	defer db.Close()

	_, err = db.Exec("DELETE FROM recover")

	if err != nil {

		log.Println(err)

		return
	}

	_, err = db.Exec("DELETE FROM athletes_times")

	if err != nil {

		log.Println(err)

		return
	}

	_, err = db.Exec("DELETE FROM resultados_chegada")

	if err != nil {

		log.Println(err)

		return
	}

	_, err = db.Exec("DELETE FROM resultados_largada")

	if err != nil {

		log.Println(err)

		return
	}

	_, err = db.Exec("DELETE FROM invalidos")

	if err != nil {

		log.Println(err)
	}

	return
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
