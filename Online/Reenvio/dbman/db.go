package dbman

import (
	"errors"
	"fmt"
	"log"
	"os"

	"database/sql"
	"github.com/MyTempoESP/Reenvio/atleta"
	//"sync/atomic"

	_ "modernc.org/sqlite"
)

type Baselet struct {
	Path string

	db     *sql.DB
	opened bool

	Tempos <-chan atleta.Atleta
}

type MADB struct { // client version
	DatabaseRoot string // path to the database dir
	databases    []Baselet
}

func NewBaselet(path string) (b Baselet, err error) {

	b.Path = path

	err = b.Init()

	return
}

func (b *Baselet) Init() (err error) {

	err = b.Open()

	if err != nil {

		return
	}

	b.beginMonitor()

	return
}

func (b *Baselet) Open() (err error) {

	if b.opened {

		return
	}

	if _, err = os.Stat(b.Path); errors.Is(err, os.ErrNotExist) {

		log.Printf("Arquivo inexistente: %s\n", b.Path)

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

func (b *Baselet) ScanCheckpoint(query string, tempos chan<- atleta.Atleta) {

	res, err := b.db.Query(query)

	if err != nil {

		log.Println("Erro checando atletas disponíveis", err)

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

			log.Println("Erro ao escanear os atletas: ", err)

			break
		}

		tempos <- at
	}

	err = res.Err()

	if err != nil {

		log.Println("Erro ao escanear os atletas: ", err)
	}
}

func (b *Baselet) Monitor() (tempos <-chan atleta.Atleta) {

	t := make(chan atleta.Atleta, 20) // groupSize

	// sqlite allows concurrent reads

	go func() {

		defer func() { close(t) }()

		b.db.Exec(ATTACH)

		b.ScanCheckpoint(QUERY_LARGADA, t)
		b.ScanCheckpoint(QUERY_CHEGADA, t)

		return
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
}

func (m *MADB) Get() (lotes <-chan []atleta.Atleta) {

	l := make(chan []atleta.Atleta)

	go func() {
		defer func() { close(l) }()

		for _, b := range m.databases {

			log.Println("Recebendo lote...")

			l <- b.Get()
		}
	}()

	lotes = l

	return
}

func (m *MADB) Add() (err error) {

	var (
		b Baselet
	)

	b, err = NewBaselet(
		fmt.Sprintf("%s/N%d.db", m.DatabaseRoot, len(m.databases)))

	if err != nil {

		return
	}

	m.databases = append(m.databases, b)

	return
}

func (m *MADB) Grow(amount int) (err error) {

	for range amount {

		err = m.Add()

		if err != nil {

			return
		}
	}

	return
}

func (m *MADB) Close() {

	for _, b := range m.databases {

		log.Println("closed!")

		b.Close()
	}
}
