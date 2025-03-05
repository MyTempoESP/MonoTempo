package main

import (
	"envio/athlete"
)

func (envio *Envio) Process() {

	tags := envio.Tags

	for {

		t := <-tags

		if t.Antena == 0 {
			/*
				Antena 0 não existe
			*/

			continue
		}

		at := athlete.Atleta{
			Antena: t.Antena,         /* Antena    */
			Numero: t.Epc,            /* Numero    */
			Staff:  t.Staff,          /* Staff     */
			Tempo:  t.TempoFormatado, /* Tempo     */
			Check:  0,                /* Check     */ //TODO
		}

		/*
			Tentar salvar o atleta no banco de dados.
		*/

		envio.SalvarAtleta(at)
	}
}
