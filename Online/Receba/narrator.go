package main

import (
	"strconv"

	"go.uber.org/zap"
)

func (r *Receba) AtualizaNarrator(path string, atletas []Atleta, logger *zap.Logger) {
	w, f, err := CreateCSVWriter(path)

	logger = logger.With(zap.String("npath", path))

	if err != nil {
		logger.Error("Error creating csv", zap.Error(err))
		return
	}
	defer f.Close()

	for _, at := range atletas {
		logger := logger.With(zap.String("Atleta", at.Nome), zap.Int("Numero", at.Numero))
		WriteCSVRecord(w, []string{strconv.Itoa(at.Numero), at.Nome}, logger)
	}
	w.Flush()
}
