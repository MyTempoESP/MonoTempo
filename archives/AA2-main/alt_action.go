package main

import (
	"log"

	mysql "github.com/mytempoesp/mysql-easy"
)

func RestartComputer() (err error) {

	db, err := mysql.ConfiguraDB()

	if err != nil {

		return
	}

	_, err = db.Exec("UPDATE equipamento SET action = 1")

	log.Println(err)

	defer db.Close()

	return
}
