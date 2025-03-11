package main

import (
	"encoding/json"

	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	backoff "github.com/cenkalti/backoff"
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
	Action  int    `json:"action"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Form map[string]string

/*
	TODO: this request module is getting too big, time
	to turn it into a repo/gist or own project.

	By Rodrigo Monteiro Junior

	- Version 0.1:

		receives a 'data' parameter
		as a map[string]string.

	- Version 0.1.1:
	  ter 10 set 2024 14:24:16 -03

		Minor patch to check errors
		related to the `status` field.

	- Version 0.2 * ( FIXME: CONFLICTING patch with 'github.com/mytempoesp/Envio/request.go' ):
	  ter 10 set 2024 15:07:49 -03

		Major patch to error reporting, won't affect
		much usage, but conform to proper idiomatic
		error handling and avoid redundancy.

Execute an HTTP POST request to a JSON api,
passing a JSON form and getting a response in
a user-defined struct.

this function retries to make the request up to 20 seconds
using a backoff algorithm.
*/
func JSONRequest(url string, data Form, jsonOutput interface{}) (err error) {

	var res *http.Response

	jsonData, err := json.Marshal(data)

	if err != nil {
		err = fmt.Errorf("Error marshaling JSON: %s\n", err)

		return
	}

	bf := backoff.NewExponentialBackOff()
	bf.MaxElapsedTime = 20 * time.Second

	err = backoff.Retry(
		func() (err error) {

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

			if err != nil {
				err = fmt.Errorf("Error creating request: %s\n", err)

				return
			}

			req.Header.Set("Content-Type", "application/json")

			res, err = http.DefaultClient.Do(req)

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
		err = fmt.Errorf("Error connecting to '%s': got HTTP %d", url, res.StatusCode)

		return
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		err = fmt.Errorf("Error reading response body: %s\n", err)

		return
	}

	//log.Println(string(body)) // XXX: Debugging

	/*
		By Rodrigo Monteiro Junior
		ter 10 set 2024 14:30:47 -03

		patch for checking a `status` response.
		(this is a nasty workaround for faster debugging)
	*/
	var check RespostaAPI

	err = json.Unmarshal(body, &check)

	if err != nil {
		log.Printf("WARN: Can't unmarshal response JSON into type %T, %s\n", check, err)

		/* we can safely ignore this, since it's simply meant for error reporting */
	}

	if check.Status == "error" {
		err = fmt.Errorf("API returned error status: %s\n", check.Message)

		return
	}
	/*
		patch is over
	*/

	err = json.Unmarshal(body, &jsonOutput)

	if err != nil {
		err = fmt.Errorf("Error unmarshaling response JSON: %s\n", err)
	}

	return
}

func GetAction(url string, data Form) (action int, err error) {

	var res *http.Response

	jsonData, err := json.Marshal(data)

	if err != nil {
		err = fmt.Errorf("Error marshaling JSON: %s\n", err)

		return
	}

	bf := backoff.NewExponentialBackOff()
	bf.MaxElapsedTime = 20 * time.Second

	err = backoff.Retry(
		func() (err error) {

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

			if err != nil {
				err = fmt.Errorf("Error creating request: %s\n", err)

				return
			}

			req.Header.Set("Content-Type", "application/json")

			res, err = http.DefaultClient.Do(req)

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
		err = fmt.Errorf("Error connecting to '%s': got HTTP %d", url, res.StatusCode)

		return
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		err = fmt.Errorf("Error reading response body: %s\n", err)

		return
	}

	//log.Println(string(body)) // XXX: Debugging

	var check RespostaAPI

	err = json.Unmarshal(body, &check)

	if err != nil {
		err = fmt.Errorf("Error: Can't unmarshal response JSON into type %T. %s\n", check, err)

		return
	}

	if check.Status == "error" {
		err = fmt.Errorf("API returned error status: %s\n", check.Message)

		return
	}

	action = check.Action

	return
}

/*
Do a simple POST request and ignore the response,
only treat it in case of errors etc.
*/
func JSONSimpleRequest(url string, data Form) (err error) {

	var res *http.Response

	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}

	bf := backoff.NewExponentialBackOff()
	bf.MaxElapsedTime = 20 * time.Second

	err = backoff.Retry(

		func() (err error) {
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

			if err != nil {
				log.Println("Error creating request:", err)

				return
			}

			req.Header.Set("Content-Type", "application/json")

			res, err = http.DefaultClient.Do(req)

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
		err = fmt.Errorf("Error connecting to '%s': got HTTP %d", url, res.StatusCode)
	}

	return
}
