package main

import (
	"log"
	"os"
	"time"

	backoff "github.com/cenkalti/backoff"
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

	err := r.Equip.Atualiza()

	if err != nil {

		log.Fatalf("Erro atualizando equipamento: %s\n", err)
	}

	r.Tempos.DatabaseRoot = "/var/monotempo-data/"

	n, err := countDir(r.Tempos.DatabaseRoot)

	if err != nil {

		log.Fatalf("Couldn't count files: %s\n", err)
	}

	log.Printf("Processing %d databases...\n", n)

	r.Tempos.Grow(n)

	lotes := r.Tempos.Get()

	if err != nil {

		log.Fatalf("Couldn't Get data: %s\n", err)
	}

	bf := backoff.NewExponentialBackOff()

	bf.MaxElapsedTime = 5 * time.Minute
	bf.MaxInterval = 10 * time.Second

	err = backoff.Retry(
		func() (err error) {

			err = r.TentarReenvio(lotes)

			if err != nil {

				log.Println(err)
			}

			return
		},

		bf,
	)

	if err != nil {

		log.Printf("Não foi possível enviar dados: %s\n", err)

		os.Exit(1)
	}
}
