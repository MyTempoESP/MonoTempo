package pinger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

type Form map[string]string

func JSONSimpleRequest(url string, data Form) (err error) {

	var (
		res      *http.Response
		req      *http.Request
		jsonData []byte
	)

	jsonData, err = json.Marshal(data)

	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}

	req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

	if err != nil {
		log.Println("Error creating request:", err)

		return
	}

	req.Header.Set("Content-Type", "application/json")

	res, err = http.DefaultClient.Do(req)

	if err != nil {
		log.Println("Error sending request:", err)

		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("Error connecting to '%s': got HTTP %d", url, res.StatusCode)
	}

	return
}

func NewJSONPinger(state *atomic.Bool) {

	url := os.Getenv("MYTEMPO_API_URL")
	infoRota := fmt.Sprintf("http://%s/status/device", url)

	tick := time.NewTicker(4 * time.Second)
	data := Form{
		"deviceId": os.Getenv("MYTEMPO_EQUIP"),
	}

	for {
		<-tick.C

		state.Store(JSONSimpleRequest(infoRota, data) == nil)
	}
}
