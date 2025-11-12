package pinger

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
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

func NewJSONPinger(state *atomic.Bool, logger *zap.Logger) {

	url := os.Getenv("MYTEMPO_API_URL")
	infoRota := fmt.Sprintf("http://%s/status/device", url)
	devRota := fmt.Sprintf("http://%s/fetch/device", url)

	devid, fetchErr := BuscaID(devRota)

	tick := time.NewTicker(14 * time.Second)

	data := Form{
		"deviceId": devid,
	}

	logger = logger.With(
		zap.String("Base URL", url),
		zap.String("Info URL", infoRota),
		zap.String("Dev URL", devRota),
	)

	for {
		<-tick.C

		if fetchErr != nil {
			devid, fetchErr = BuscaID(devRota)

			data = Form{
				"deviceId": devid,
			}

			logger = logger.With(
				zap.String("Device ID", devid),
			)
		}

		logger.Info("Sending JSON request to INFO URL")

		err := JSONSimpleRequest(infoRota, data)

		logger.Info("Request terminated")

		state.Store(err == nil)

		if err != nil {
			logger.Error("Request error", zap.Error(err))
		}
	}
}
