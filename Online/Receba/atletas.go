package main

import (
	//"log"
	"strconv"
)

/*
{
	"numero": 420,
	"provaId": 201,
	"percursoId": 401,
	"nome": "DALVA",
	"equipe": "JATO"
}
*/

type Atleta struct {
	Nome       string `json:"nome"`
	Sexo       string `json:"sexo"`
	Equipe     string `json:"equipe"`
	Cidade     string `json:"cidade"`
	Numero     int    `json:"numero"`
	ProvaID    int    `json:"provaId"`
	PercursoID int    `json:"percursoId"`
}

func (r *Receba) BuscaAtletas(idProva int) (atletas []Atleta, err error) {

	data := Form{
		"idProva": strconv.Itoa(idProva),
	}

	err = JSONRequest(r.AtletasRota, data, &atletas)

	return
}

func (r *Receba) AtualizaAtletas(atletas []Atleta) (err error) {

	for _, a := range atletas {
		_, err = r.db.Exec(
			QUERY_ATUALIZA_ATLETA,

			a.Nome,
			a.Sexo,
			a.Equipe,
			a.Cidade, /* city */
			a.Numero,
			a.ProvaID,
			a.PercursoID,
		)

		//log.Println(a.Sexo)
	}

	return
}
