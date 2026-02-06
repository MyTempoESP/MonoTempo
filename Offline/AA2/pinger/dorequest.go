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
	security := os.Getenv("API_REQ_SECURITY")

	infoRota := fmt.Sprintf("%s://%s/status/device", security, url)
	devRota := fmt.Sprintf("%s://%s/fetch/device", security, url)

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

			logger.Info("Got device id",
				zap.String("Device ID", devid),
			)
		}

		logger_new := logger.With(zap.String("device id", devid))
		logger_new.Info("Sending JSON request to INFO URL")

		err := JSONSimpleRequest(infoRota, data)

		logger_new.Info("Request terminated")

		state.Store(err == nil)

		if err != nil {
			logger_new.Error("Request error", zap.Error(err))
		}
	}
}