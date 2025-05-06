package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	/* i love this library */
	backoff "github.com/cenkalti/backoff"
	"go.uber.org/zap"
)

/*
By Rodrigo Monteiro Junior
ter 10 set 2024 14:24:16 -03

-- FROM V0.2 --

Resposta genérica da API do kerlo
contendo apenas fields relacionados
a status e mensagens de sucesso/falha.
*/
type RespostaAPI struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var (
	ErrNetwork  = errors.New("erro de rede")
	ErrBodyRead = errors.New("erro lendo body")
)

type APIError struct {
	Message string
}

func (e *APIError) Error() string {
	return e.Message
}

const (
	REQUEST_TIMEOUT = 10 * time.Second
)

type Form map[string]string
type RawForm []byte

func SimpleRawRequest(url string, data RawForm, contentType string, logger *zap.Logger) (err error) {

	var res *http.Response

	bf := backoff.NewExponentialBackOff()
	bf.MaxElapsedTime = REQUEST_TIMEOUT

	err = backoff.Retry(
		func() (err error) {
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))

			if err != nil {
				return
			}

			req.Header.Set("Content-Type", contentType)

			res, err = http.DefaultClient.Do(req)

			return
		},

		bf,
	)

	if err != nil {

		err = ErrNetwork

		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = ErrNetwork

		// err = fmt.Errorf("error connecting to '%s': got HTTP %d", url, res.StatusCode)

		return
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		err = ErrBodyRead

		return
	}

	/*
		By Rodrigo Monteiro Junior
		ter 10 set 2024 14:30:47 -03

		patch for checking a `status` response.
		(this is a nasty workaround for faster debugging)
	*/
	var check RespostaAPI

	jsonErr := json.Unmarshal(body, &check)

	if jsonErr != nil {
		/* we can safely ignore this, since it's simply meant for error reporting */
		logger.Warn("Json error", zap.Error(err))
	} else {
		if check.Status == "error" {

			err = &APIError{check.Message}

			return
		}
	}

	return
}
