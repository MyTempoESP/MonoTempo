package main

import (
	"errors"

	"go.uber.org/zap"
)

type Equipamento struct {
	ID      int    `json:"id"`
	Nome    string `json:"modelo"`
	ProvaID int    `json:"assocProva"`
	Check   int    `json:"assocCheck"`
}

var ErrEquipAssoc = errors.New("equip nao assoc")

func (r *Receba) BuscaEquip(equipModelo string, logger *zap.Logger) (equip Equipamento, err error) {

	var ae *APIError

	data := Form{
		"device": equipModelo,
	}

	err = JSONRequest(r.DeviceRota, data, &equip, logger)

	if errors.Is(err, ErrNetwork) {
		return
	}

	if errors.As(err, &ae) {
		return
	}

	if equip.ProvaID == 0 {
		err = ErrEquipAssoc
	}

	return
}

func (r *Receba) AtualizaEquip(equip Equipamento) (err error) {

	_, err = r.db.Exec(
		queryAtualizaEquip,

		equip.ID,
		equip.Nome,
		equip.ProvaID,
		equip.Check,
	)

	return
}
