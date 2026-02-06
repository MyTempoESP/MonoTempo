package dbman

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"database/sql"

	"github.com/MyTempoESP/Reenvio/atleta"
	"go.uber.org/zap"

	_ "modernc.org/sqlite"
)

type Baselet struct {
	Path         string
	IsCheckpoint bool // true if this is a checkpoint database

	db     *sql.DB
	opened bool

	Logger *zap.Logger

	Tempos <-chan atleta.Atleta
}

type MADB struct { // client version
	DatabaseRoot string // path to the database dir
	IsCheckpoint bool   // true if this is a checkpoint database

	Logger *zap.Logger

	databases []Baselet
}

func NewBaselet(path string, check bool, n int, logger *zap.Logger) (b Baselet, err error) {

	b.IsCheckpoint = check
	b.Path = path

	b.Logger = logger.With(
		zap.Int("baselet_id", n),
		zap.String("db_sub_path", path),
	)

	err = b.Init()

	return
}

func (b *Baselet) Init() (err error) {

	b.Logger.Info("Iniciando baselet...")

	err = b.Open()

	if err != nil {

		b.Logger.Error("Erro ao abrir baselet", zap.Error(err))
		return
	}

	b.Logger.Info("Baselet aberto com sucesso!")

	b.beginMonitor()

	return
}

func (b *Baselet) Open() (err error) {

	if b.opened {

		return
	}

	if _, err = os.Stat(b.Path); errors.Is(err, os.ErrNotExist) {

		return
	}

	db, err := sql.Open("sqlite", b.Path)

	if err != nil {

		return
	}

	b.db = db

	b.opened = true

	return
}

func (b *Baselet) beginMonitor() {

	// create data channel for insertions
	b.Tempos = b.Monitor()
}

func (b *Baselet) ScanCheckpoint(query string, logger *zap.Logger, tempos chan<- atleta.Atleta) {

	startTime := time.Now()

	logger.Info("Iniciando leitura de checkpoint...")

	res, err := b.db.Query(query)

	if err != nil {

		logger.Error("Erro ao escanear os atletas",
			zap.Error(err),
			zap.Duration("duration", time.Since(startTime)),
		)

		return
	}

	defer res.Close()

	for res.Next() {

		var at atleta.Atleta

		err = res.Scan(
			&at.Numero,
			&at.Antena,
			&at.PercursoID,
			&at.Tempo,
		)

		if err != nil {

			logger.Error("Erro ao escanear os atletas obtidos na query",
				zap.Error(err),
				zap.Duration("duration", time.Since(startTime)),
			)

			return
		}

		tempos <- at
	}

	err = res.Err()

	if err != nil {

		logger.Error("Erro ao escanear os atletas obtidos na query",
			zap.Error(err),
			zap.Duration("duration", time.Since(startTime)),
		)

		return
	}

	logger.Info("Leitura de checkpoint finalizada!",
		zap.Duration("duration", time.Since(startTime)),
	)
}

func (b *Baselet) Monitor() (tempos <-chan atleta.Atleta) {

	b.Logger.Info("Iniciando monitoramento...")

	t := make(chan atleta.Atleta, 20) // groupSize

	// sqlite allows concurrent reads

	go func() {

		defer func() { close(t) }()
		defer b.Logger.Info("Monitoramento encerrado!")

		_, err := b.db.Exec(ATTACH)

		if err != nil {
			b.Logger.Error("Erro obtendo dados do equipamento", zap.Error(err))
			return
		}

		if !b.IsCheckpoint {
			logger := b.Logger.With(zap.String("checkpoint_type", "largada"))
			b.ScanCheckpoint(QUERY_LARGADA, logger, t)

			logger = b.Logger.With(zap.String("checkpoint_type", "chegada"))
			b.ScanCheckpoint(QUERY_CHEGADA, logger, t)
		} else {
			logger := b.Logger.With(zap.String("checkpoint_type", "checkpoint"))
			b.ScanCheckpoint(QUERY_CHECKPOINT, logger, t)
		}
	}()

	tempos = t

	return
}

func (b *Baselet) Get() (atletas []atleta.Atleta) {

	var data atleta.Atleta

	for data = range b.Tempos {
		atletas = append(atletas, data)
	}

	return
}

func (b *Baselet) Close() {

	if !b.opened {

		return
	}

	b.db.Close()
	b.opened = false

	b.Logger.Info("Baselet fechado com sucesso!")
}

func (m *MADB) Get() (lotes <-chan []atleta.Atleta) {

	l := make(chan []atleta.Atleta)

	go func() {
		var wg sync.WaitGroup

		defer func() {
			wg.Wait()
			close(l)
		}()

		for n, b := range m.databases {
			wg.Add(1)

			go func(n int, b Baselet) {
				defer wg.Done()

				logger := m.Logger.With(
					zap.Int("baselet_id", n))

				startTime := time.Now()
				logger.Info("Recebendo lote...")

				a := b.Get()
				l <- a

				logger.Info("Lote recebido com sucesso",
					zap.Duration("duration", time.Since(startTime)),
					zap.Int("batch_size", len(a)),
				)
			}(n, b)
		}
	}()

	lotes = l

	return
}

func (m *MADB) Add() (err error) {

	var (
		b Baselet
	)

	n := len(m.databases)

	b, err = NewBaselet(
		fmt.Sprintf("%s/N%d.db", m.DatabaseRoot, n),
		m.IsCheckpoint,
		n,
		m.Logger,
	)

	if err != nil {

		return
	}

	m.databases = append(m.databases, b)

	return
}

func (m *MADB) Grow(amount int) (err error) {

	m.Logger.Info("Crescendo MADB...",
		zap.Int("amount", amount))

	if amount < 1 {
		return
	}

	for range amount {

		err = m.Add()

		if err != nil {

			return
		}
	}

	return
}

func (m *MADB) Close() {

	m.Logger.Info("Fechando MADB...")

	for _, b := range m.databases {

		b.Close()
	}
}
