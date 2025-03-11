package main

import (
	"envio/athlete"
)

func (envio *Envio) SalvarAtleta(a athlete.Atleta) {
	envio.DBManager.Insert(a)
}
