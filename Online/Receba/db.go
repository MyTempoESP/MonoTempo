package main

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

func (r *Receba) ConfiguraDB() (err error) {

	db, err := sql.Open("sqlite", "/var/monotempo-data/equipamento.db")

	if err != nil {

		return
	}

	_, err = db.Exec(CRIA_DB)

	if err != nil {

		return
	}

	r.db = db

	return
}
