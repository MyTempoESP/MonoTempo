package main

import (
	"github.com/MyTempoESP/Validate/dbman"
	"go.uber.org/zap"
)

type Reenvio struct {
	Tempos dbman.MADB
	Equip  Equipamento
	Logger *zap.Logger
}
