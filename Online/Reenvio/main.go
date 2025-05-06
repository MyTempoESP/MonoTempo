package main

import (
	"errors"
	"os"
	"time"

	backoff "github.com/cenkalti/backoff"
	"go.uber.org/zap"
)

func countDir(path string) (n int, err error) {

	f, err := os.Open(path)

	if err != nil {

		return
	}

	list, err := f.Readdirnames(-1)

	f.Close()

	if err != nil {

		return
	}

	n = len(list)

	return
}

func main() {

	var r Reenvio

	logger, err := zap.NewProduction()

	if err != nil {
		panic(err)
	}

	r.Logger = logger

	err = r.Equip.Atualiza()

	logger = r.Logger.With(
		zap.String("equipamento", r.Equip.Nome),
		zap.Int("equip_id", r.Equip.ID),
		zap.Int("prova_id", r.Equip.ProvaID),
		zap.Int("check_id", r.Equip.Check),
	)

	logger.Info("Equipamento atualizado, iniciando reenvio de dados...")

	if err != nil {
		logger.Error("Erro ao atualizar equipamento", zap.Error(err))
		os.Exit(1)
	}

	if r.Equip.Check != 0 {
		logger.Debug("Checkpoint detectado, modo automático desativado")
		r.Tempos.IsCheckpoint = true
	}

	r.Tempos.DatabaseRoot = "/var/monotempo-data/"

	logger = r.Logger.With(
		zap.String("database_root", r.Tempos.DatabaseRoot),
	)

	logger.Info("Contando arquivos de banco de dados")

	n, err := countDir(r.Tempos.DatabaseRoot)

	if err != nil {
		logger.Error("Erro ao contar arquivos de banco de dados", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("Arquivos encontrados, iniciando MADB",
		zap.Int("databases", n))

	r.Tempos.Logger = logger.With(
		zap.String("db_root_path", r.Tempos.DatabaseRoot),
		zap.Bool("is_checkpoint", r.Tempos.IsCheckpoint),
	)

	r.Tempos.Grow(n - 1)

	lotes := r.Tempos.Get()

	r.Logger.Info("Enviando dados")

	bf := backoff.NewExponentialBackOff()

	bf.MaxElapsedTime = 5 * time.Minute
	bf.MaxInterval = 10 * time.Second

	err = backoff.Retry(
		func() (err error) {

			err = r.TentarReenvio(lotes, r.Logger)

			if errors.Is(err, ErrWrongDate) { // if date is wrong, don't even retry
				err = backoff.Permanent(err)
			}

			return
		},

		bf,
	)

	if err != nil {

		logger.Error("Erro no reenvio", zap.Error(err))
		os.Exit(1)
	}
}
