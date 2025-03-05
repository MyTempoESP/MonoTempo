package dbman

import (
	"fmt"
	"log"
	"os"

	"database/sql"
	"envio/athlete"
	_ "modernc.org/sqlite"
)

type Baselet struct {
	Path string

	db     *sql.DB
	data   chan<- athlete.Atleta
	opened bool
}

type MADB struct {
	DatabaseRoot string // path to the database dir

	databases []Baselet
	groupSize int // group length in which athletes are divided (0-10, 0-100, ...)
	maxValue  int // max number that fits
}

func NewBaselet(path string) (b Baselet, err error) {

	_, err = os.Create(path)

	if err != nil {

		return
	}

	db, err := sql.Open("sqlite", path)

	if err != nil {

		return
	}

	defer db.Close() // flush changes

	_, err = db.Exec(CREATE_TIME_DATABASE)

	if err != nil {

		return
	}

	b.db = db
	b.Path = path

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

	// create data channel for insertions
	b.data = b.Monitor()

	return
}

func (b *Baselet) Insert(c athlete.Atleta) (err error) {

	err = b.Open()

	if err != nil {

		return
	}

	b.data <- c

	return
}

func (b *Baselet) Monitor() chan<- athlete.Atleta {

	data := make(chan athlete.Atleta)

	go func() {
		for c := range data {

			b.db.Exec(
				INSERT_TIME,

				c.Antena,
				c.Numero,
				c.Staff,
				c.Tempo,
			)
		}
	}()

	return data
}

func (b *Baselet) Close() {

	if !b.opened {

		return
	}

	close(b.data)

	b.db.Close()

	b.opened = false
}

func (m *MADB) GroupSize(s int) {
	m.groupSize = s
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

	m.maxValue += m.groupSize * amount

	return
}

func (m *MADB) Insert(c athlete.Atleta) (err error) {

	if c.Numero > m.maxValue {

		err = m.Grow(((c.Numero - m.maxValue) / m.groupSize) + 1)

		if err != nil {

			return
		}
	}

	if c.Numero == m.maxValue {

		err = m.Grow(1)

		if err != nil {

			return
		}
	}

	err = m.databases[c.Numero/m.groupSize].Insert(c)

	return
}

func (m *MADB) Close() {

	for _, b := range m.databases {

		log.Println("closed!")

		b.Close()
	}
}

func (m *MADB) Init() (err error) {

	err = m.Grow(1)

	return
}
