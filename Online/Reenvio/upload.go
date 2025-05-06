package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"encoding/json"

	"github.com/MyTempoESP/Reenvio/atleta"
	"github.com/MyTempoESP/Reenvio/narrator"
	"go.uber.org/zap"
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

var ErrWrongDate = errors.New("a data da prova não coincide com a data atual")

func (reenvio *Reenvio) Upload(atletas []atleta.Atleta, voicelog narrator.Narrator, logger *zap.Logger) (err error) {

	startTime := time.Now()

	logger = logger.With(
		zap.Int("athlete_count", len(atletas)),
	)

	if len(atletas) == 0 {

		logger.Warn("Lista de atletas vazia")

		return
	}

	if reenvio.Equip.ID == 0 {

		logger.Error("Equip inválido")

		return
	}

	logger = logger.With(
		zap.Int("equipamento_id", reenvio.Equip.ID),
	)

	if reenvio.Equip.ProvaID == 0 {

		logger.Error("Prova inválida")

		return
	}

	logger = logger.With(
		zap.Int("prova_id", reenvio.Equip.ProvaID),
	)

	dados := AtletasForm{
		reenvio.Equip.ID,
		reenvio.Equip.ProvaID,
		atletas,
	}

	data, err := json.Marshal(dados)

	if err != nil {

		logger.Error("Erro ao serializar dados",
			zap.Error(err))

		return
	}

	/* NOTE: DEBUG */
	logger.Debug("Dados a serem enviados",
		zap.String("dados", string(data)),
	)

	logger = logger.With(zap.String("endpoint", UrlTempos))
	logger.Info("Enviando dados para a API")

	/*
		O erro na request não é retornado, apenas lidamos com ele.
	*/
	err = SimpleRawRequest(UrlTempos, data, "application/json", logger)

	if err != nil {
		logger.Error("Erro ao enviar dados",
			zap.Error(err),
			zap.Duration("tempo", time.Since(startTime)),
		)

		if errors.Is(err, ErrNetwork) {
			voicelog.SayString("Erro de internet, verifique a conexão")
		}

		var ae *APIError

		if errors.As(err, &ae) {
			logger.Warn("API Level error detected", zap.Error(err))

			if strings.ToUpper(ae.Message) == "A DATA DA PROVA NÃO COINCIDE COM A DATA ATUAL." {
				err = ErrWrongDate
			}
		}

		return
	}

	logger.Info("Dados enviados com sucesso",
		zap.Duration("tempo", time.Since(startTime)),
	)

	return
}

/*
By Rodrigo Monteiro Junior
ter 17 set 2024 12:34:47 -03

# TentarReenvio

Aguarda por dados de reenvio e tenta o envio para a API,
em caso de erro, redireciona para o Banco de Dados.
*/
func (reenvio *Reenvio) TentarReenvio(lotes <-chan []atleta.Atleta, logger *zap.Logger) (err error) {

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

	logger = logger.With(zap.Duration("timeout", TIMEOUT_MONITORAMENTO))

	/*
		Se houver registros, receba.
		Caso contrário, não bloqueie o código.
	*/
	vl := narrator.New()
	vl.Enabled = true
	vl.Watch()

	for {
		select {
		case tempos, ok = <-lotes:
			if !ok {
				return
			}

			uploadErr := reenvio.Upload(tempos, vl, logger)

			if errors.Is(uploadErr, ErrWrongDate) {
				vl.SayString(ErrWrongDate.Error())
				err = uploadErr
			}

		case <-timeoutMon:
			logger.Warn("Timeout de monitoramento atingido")
			err = fmt.Errorf("timeout, deixando para enviar depois")
			return
		}
	}
}
