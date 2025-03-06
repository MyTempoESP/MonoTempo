package main

import (
	"database/sql"
)

const (
	QUERY_ATUALIZA_EQUIP = `
	REPLACE INTO equipamento (
	    id,
	    idequip,
	    modelo,
	    event_id
	)
	VALUES(1, ?, ?, ?);
	`

	QUERY_ATUALIZA_PROVA = `
	REPLACE INTO event_data (
		id,
		event_date,
		event_title
	)
	VALUES (?, ?, ?) 
	`

	QUERY_ATUALIZA_PERCURSO = `
	REPLACE INTO tracks (
	    id,
	    race_description,
	    event_id,
	    inicio,
	    chegada,
	    largada
	)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	QUERY_ATUALIZA_ATLETA = `
	REPLACE INTO athletes (
	    name,
	    sex,
	    team,
	    city,
	    num,
	    event_id,
	    track_id
	)
	VALUES(?, ?, ?, ?, ?, ?, ?);
	`

	QUERY_ATUALIZA_STAFF = `
	REPLACE INTO staffs (
	    id,
	    event_id,
	    nome
	)
	VALUES(?, ?, ?);
	`
)

func IgnorarForeignKey(db *sql.DB) {
	db.Exec("SET FOREIGN_KEY_CHECKS=0")
}

func AceitarForeignKey(db *sql.DB) {
	db.Exec("SET FOREIGN_KEY_CHECKS=1")
}
