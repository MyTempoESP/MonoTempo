package pinger

import (
	"encoding/json"
	"time"

	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
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
*/
func JSONRequest(url string, data Form, jsonOutput interface{}) (err error) {

	var res *http.Response

	jsonData, err := json.Marshal(data)

	if err != nil {
		err = fmt.Errorf("error marshaling JSON: %s", err)

		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	if err != nil {
		err = fmt.Errorf("error creating request: %s", err)

		return
	}

	req.Header.Set("Content-Type", "application/json")

	res, err = http.DefaultClient.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("error connecting to '%s': got HTTP %d", url, res.StatusCode)

		return
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		err = fmt.Errorf("error reading response body: %s", err)

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
		err = fmt.Errorf("api returned error status: %s", check.Message)

		return
	}
	/*
		patch is over
	*/

	err = json.Unmarshal(body, &jsonOutput)

	if err != nil {
		err = fmt.Errorf("error unmarshaling response JSON: %s", err)
	}

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

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	if err != nil {
		log.Println("Error creating request:", err)

		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 15 * time.Second, // Set your desired timeout here
	}

	res, err = client.Do(req)

	if err != nil {
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("error connecting to '%s': got HTTP %d", url, res.StatusCode)
	}

	return
}
