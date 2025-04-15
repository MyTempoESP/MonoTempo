package pinger

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"
)

type Equipamento struct {
	ID      int    `json:"id"`
	Nome    string `json:"modelo"`
	ProvaID int    `json:"assocProva"`
}

func BuscaEquip(equipModelo string, url string) (equip Equipamento, err error) {

	data := Form{
		"device": equipModelo,
	}

	err = JSONRequest(url, data, &equip)

	return
}

func BuscaID(url string) (devid string, err error) {

	equip, err := BuscaEquip(os.Getenv("MYTEMPO_EQUIP"), url)

	devid = "0"

	if err != nil {
		log.Println("Error fetching device, won't comm", err)
	} else {
		devid = fmt.Sprintf("%d", equip.ID)
		log.Println("Device ID:", devid)
	}

	return
}

func NewJSONPinger(state *atomic.Bool) {

	url := os.Getenv("MYTEMPO_API_URL")
	infoRota := fmt.Sprintf("http://%s/status/device", url)
	devRota := fmt.Sprintf("http://%s/fetch/device", url)

	devid, fetchErr := BuscaID(devRota)

	tick := time.NewTicker(4 * time.Second)

	data := Form{
		"deviceId": devid,
	}

	for {
		<-tick.C

		if fetchErr != nil {
			devid, fetchErr = BuscaID(devRota)

			data = Form{
				"deviceId": devid,
			}
		}

		log.Println("Sending JSON request to", infoRota)

		err := JSONSimpleRequest(infoRota, data)

		log.Println("Request terminated")

		state.Store(err == nil)
		log.Println(err)
	}
}
