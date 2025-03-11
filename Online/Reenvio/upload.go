package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"encoding/json"
	"github.com/MyTempoESP/Reenvio/atleta"
)

const (
	REENVIO_INTERVALO     = 23 * time.Second
	TIMEOUT_MONITORAMENTO = 5 * time.Second
)

var (
	/*
		URL para subir os tempos.
	*/
	UrlTempos = fmt.Sprintf(
		"http://%s/receive/tempos", os.Getenv("MYTEMPO_API_URL"))
)

/*
	{
	    "equipamentoId": 103,
	    "provaId": 201,
	    "atletas": [
	    	{
			"antena": 1,
			"numero": 201,
			"tempo": "10:20:33.323",
			"percurso": 120
		}
	    ]
	}
*/
type AtletasForm struct {
	EquipamentoId int             `json:"equipamentoId"`
	ProvaId       int             `json:"provaId"`
	Atletas       []atleta.Atleta `json:"atletas"`
}

func (reenvio *Reenvio) Upload(atletas []atleta.Atleta) (err error) {

	if len(atletas) == 0 {

		log.Println("Ignorando lote vazio de atletas para reenvio.")

		return
	}

	if reenvio.Equip.ID == 0 {

		err = fmt.Errorf("Equipamento inválido! (0)\n")

		return
	}

	if reenvio.Equip.ProvaID == 0 {

		err = fmt.Errorf("Prova inválida! (0)\n")

		return
	}

	dados := AtletasForm{
		reenvio.Equip.ID,
		reenvio.Equip.ProvaID,
		atletas,
	}

	data, err := json.Marshal(dados)

	if err != nil {

		err = fmt.Errorf("Erro na conversão dos dados: %s\n", err)

		return
	}

	/* NOTE: DEBUG */
	log.Printf("Dados para a request: %s\n", string(data))

	/*
		O erro na request não é retornado, apenas lidamos com ele.
	*/
	uploadErr := SimpleRawRequest(UrlTempos, data, "application/json")

	if uploadErr != nil {

		err = fmt.Errorf("Tentativa de reenvio de um lote falhou. (erro: %s)\n", uploadErr)
	}

	return
}

/*
By Rodrigo Monteiro Junior
ter 17 set 2024 12:34:47 -03

# TentarReenvio

Aguarda por dados de reenvio e tenta o envio para a API,
em caso de erro, redireciona para o Banco de Dados.
*/
func (reenvio *Reenvio) TentarReenvio(lotes <-chan []atleta.Atleta) (err error) {

	var (
		tempos []atleta.Atleta
		ok     bool
	)

	/*
		Envia um lote.
	*/

	/*
		Criar o timeout para o monitoramento.
		define o tempo gasto esperando um tempo vir da lista.
	*/
	timeoutMon := time.After(
		min(TIMEOUT_MONITORAMENTO, REENVIO_INTERVALO-1),
	)

	/*
		Se houver registros, receba.
		Caso contrário, não bloqueie o código.
	*/
	for {
		select {
		case tempos, ok = <-lotes:
			if !ok {

				return
			}

			err = reenvio.Upload(tempos)

			if err != nil {
				return
			}
		case <-timeoutMon:
			err = fmt.Errorf("Timeout, deixando para enviar depois")
			return
		}
	}
}
