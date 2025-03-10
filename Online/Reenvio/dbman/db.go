package dbman

import (
	"fmt"
	"log"

	"database/sql"
	"github.com/MyTempoESP/Reenvio/atleta"
	//"sync/atomic"

	_ "modernc.org/sqlite"
)

type Baselet struct {
	Path string

	db     *sql.DB
	opened bool

	Chegadas, Largadas <-chan atleta.Atleta
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

	err = b.loadEnvironment()

	if err != nil {

		log.Printf("Error loading track data: %s\n", err)

		return
	}

	b.beginMonitor()

	return
}

func (b *Baselet) loadEnvironment() (err error) {

	_, err = b.db.Exec(ATTACH)

	return
}

func (b *Baselet) Open() (err error) {

	if b.opened {

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
	b.Largadas, b.Chegadas = b.Monitor()
}

func (b *Baselet) Monitor() (largada, chegada <-chan atleta.Atleta) {

	l := make(chan atleta.Atleta)
	c := make(chan atleta.Atleta)

	// sqlite allows concurrent reads

	go func() {

		defer func() { close(l) }()

		res, err := b.db.Query(QUERY_LARGADA)

		if err != nil {

			log.Println("Erro checando atletas disponíveis para largada", err)

			return
		}

		defer res.Close()

		for res.Next() {

			var at atleta.Atleta

			err = res.Scan(
				&at.Numero,
				&at.Staff,
				&at.Antena,
				&at.Tempo,
			)

			if err != nil {

				log.Println("Erro ao escanear os atletas: ", err)

				break
			}

			l <- at
		}

		err = res.Err()

		if err != nil {

			log.Println("Erro ao escanear os atletas: ", err)
		}

		return
	}()

	go func() {

		defer func() { close(c) }()

		res, err := b.db.Query(QUERY_CHEGADA)

		if err != nil {

			log.Println("Erro checando atletas disponíveis para chegada", err)

			return
		}

		defer res.Close()

		for res.Next() {

			var at atleta.Atleta

			err = res.Scan(
				&at.Numero,
				&at.Tempo,
				&at.Antena,
				&at.Staff,
			)

			if err != nil {

				log.Printf("Erro ao escanear atletas: %s", err)

				break
			}

			c <- at
		}

		err = res.Err()

		if err != nil {

			log.Println("Erro ao escanear os atletas: ", err)
		}

		return
	}()

	largada = l
	chegada = c

	return
}

func (b *Baselet) Get() (atletas []atleta.Atleta, err error) {

	var data atleta.Atleta
	var largada_ok, chegada_ok bool

	for {
		select {
		case data, largada_ok = <-b.Largadas:
			if !largada_ok {
				b.Largadas = nil
			}
		case data, chegada_ok = <-b.Chegadas:
			if !chegada_ok {
				b.Chegadas = nil
			}
		}

		if b.Largadas == nil && b.Chegadas == nil {

			break
		}

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
			lote, err := b.Get()

			if err != nil {

				log.Printf("Erro ao receber lote: %s\n", err)

				continue
			}

			l <- lote
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
