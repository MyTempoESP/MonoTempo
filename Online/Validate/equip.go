package main

import (
	"fmt"

	"database/sql"

	_ "modernc.org/sqlite"
)

type Equipamento struct {
	ID      int
	Nome    string
	ProvaID int
	Check   int
}

func (equip *Equipamento) Atualiza() (err error) {

	equipDB, err := sql.Open("sqlite", "/var/monotempo-data/equipamento.db")

	if err != nil {

		return
	}

	defer equipDB.Close()

	query := `SELECT idequip, modelo, event_id, check_id FROM equipamento WHERE 1;`

	res, err := equipDB.Query(query)

	if err != nil {

		return
	}

	defer res.Close()

	if !res.Next() {

		err = fmt.Errorf("dados do dispositivo não encontrados")

		return
	}

	err = res.Scan(
		&equip.ID,
		&equip.Nome,
		&equip.ProvaID,
		&equip.Check,
	)

	if err != nil {

		return
	}

	err = res.Err()

	return
}
