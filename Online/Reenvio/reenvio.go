package main

import (
	"github.com/MyTempoESP/Reenvio/dbman"
)

type Reenvio struct {
	Tempos dbman.MADB
	Equip  Equipamento
}
