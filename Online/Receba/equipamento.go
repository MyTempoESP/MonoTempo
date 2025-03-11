package main

type Equipamento struct {
	ID      int    `json:"id"`
	Nome    string `json:"modelo"`
	ProvaID int    `json:"assocProva"`
}

func (r *Receba) BuscaEquip(equipModelo string) (equip Equipamento, err error) {

	data := Form{
		"device": equipModelo,
	}

	err = JSONRequest(r.DeviceRota, data, &equip)

	return
}

func (r *Receba) AtualizaEquip(equip Equipamento) (err error) {

	_, err = r.db.Exec(
		QUERY_ATUALIZA_EQUIP,

		equip.ID,
		equip.Nome,
		equip.ProvaID,

		equip.ID,
		equip.Nome,
		equip.ProvaID,
	)

	return
}
