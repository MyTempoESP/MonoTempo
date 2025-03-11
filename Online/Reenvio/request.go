package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	/* i love this library */
	backoff "github.com/cenkalti/backoff"
)

const (
	REQUEST_TIMEOUT = 20 * time.Second
)

type Form map[string]string
type RawForm []byte

func SimpleRawRequest(url string, data RawForm, contentType string) (err error) {

	var res *http.Response

	bf := backoff.NewExponentialBackOff()
	bf.MaxElapsedTime = REQUEST_TIMEOUT

	err = backoff.Retry(
		func() (err error) {
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))

			if err != nil {
				err = fmt.Errorf("Error creating request: %s\n", err)

				return
			}

			req.Header.Set("Content-Type", contentType)

			res, err = http.DefaultClient.Do(req)

			/* FIXME: remove excessive loggin */
			if err != nil {
				log.Println("Error sending request:", err)
			}

			return
		},

		bf,
	)

	if err != nil {
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("Error connecting to '%s': got HTTP %d\n", url, res.StatusCode)

		return
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		err = fmt.Errorf("Error reading response body: %s\n", err)

		return
	}

	// NOTE: You can comment out this section safely
	log.Println("\033[31;1m " + string(body) + " \033[0m")

	return
}
