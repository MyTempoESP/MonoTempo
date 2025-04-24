package logparse

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type EquipStatus struct {
	Status      bool
	Errcount    int
	Databases   int
	UploadCount int
	AvgProctime int
}

func ParseJSONLog(filePath string) (st EquipStatus, err error) {

	f, err := os.Open(filePath)

	if err != nil {

		return
	}

	defer f.Close()

	var (
		dbTotal, dbProc, errs, total, athletes, avg int
	)

	s := bufio.NewScanner(f)
	for s.Scan() {
		var m map[string]any

		line := s.Text()

		err = json.Unmarshal([]byte(line), &m)

		if err != nil {

			continue
		}

		total++

		if m["level"] == "info" {
			switch m["msg"] {
			case "No data":
				err = fmt.Errorf("no log data")
				return
			case "Dados enviados com sucesso":
				dbProc++
				avg += int(m["duration"].(float64) * 1000)
				athletes += int(m["athlete_count"].(float64))
			case "Arquivos encontrados, iniciando MADB":
				dbTotal = int(m["databases"].(float64))
			}
		}

		if m["level"] == "error" {
			errs++
		}
	}

	st.Status = true

	st.AvgProctime = avg / dbProc
	st.Databases = dbTotal
	st.UploadCount = athletes
	st.Errcount = errs

	if dbTotal == 0 {
		st.Status = false
	}

	if st.UploadCount == 0 {
		st.Status = false
	}

	if dbProc == 0 {
		st.Status = false
	}

	if errs > (total / 2) {
		st.Status = false
	}

	return
}
