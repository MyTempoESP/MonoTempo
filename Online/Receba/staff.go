package main

import (
	"log"
	"strconv"
)

type Staff struct {
	ID   int
	Nome string
}

func (r *Receba) BuscaStaff(idProva int) (staffs []Staff, err error) {

	data := Form{
		"idProva": strconv.Itoa(idProva),
	}

	err = JSONRequest(r.StaffRota, data, &staffs)

	return
}

func (r *Receba) AtualizaStaff(staffs []Staff, idProva int) (err error) {

	for _, staff := range staffs {
		_, err := r.db.Exec(
			QUERY_ATUALIZA_STAFF,

			staff.ID,
			idProva,
			staff.Nome,
		)

		if err != nil {
			log.Printf("erro no Sql atualizando os staffs %+v\n", err)
		}
	}

	return
}
