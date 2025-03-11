package main

import (
	"strconv"
)

/*
	{
	    "id": 201,
	    "nomeProva": "TESTE 201",
	    "dataProva": "2024-07-19",
	}
*/

type Prova struct {
	ID        int        `json:"id"`
	Nome      string     `json:"nomeProva"`
	Data      string     `json:"dataProva"`
	Percursos []Percurso `json:"percursos"`
}

func (r *Receba) BuscaProva(idProva int) (prova Prova, err error) {

	data := Form{
		"idProva": strconv.Itoa(idProva),
	}

	err = JSONRequest(r.ProvaRota, data, &prova)

	return
}

func (r *Receba) AtualizaProva(prova Prova) (err error) {

	r.db.Exec(
		QUERY_ATUALIZA_PROVA,

		prova.ID,
		prova.Data,
		prova.Nome,
	)

	// atualiza cada percurso individualmente
	for _, percurso := range prova.Percursos {

		err = r.AtualizaPercurso(percurso, prova.ID)
	}

	return
}
