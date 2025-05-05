package main

import (
	"errors"
)

type Equipamento struct {
	ID      int    `json:"id"`
	Nome    string `json:"modelo"`
	ProvaID int    `json:"assocProva"`
	Check   int    `json:"assocCheck"`
}

func (r *Receba) BuscaEquip(equipModelo string) (equip Equipamento, err error) {

	var ae *APIError

	data := Form{
		"device": equipModelo,
	}

	err = JSONRequest(r.DeviceRota, data, &equip)

	if errors.Is(err, ErrNetwork) {
		Say("Erro de rede, verifique a conexão")
	}

	if errors.As(err, &ae) {
		Say(err.Error())
	}

	return
}

func (r *Receba) AtualizaEquip(equip Equipamento) (err error) {

	_, err = r.db.Exec(
		QUERY_ATUALIZA_EQUIP,

		equip.ID,
		equip.Nome,
		equip.ProvaID,
		equip.Check,
	)

	return
}
