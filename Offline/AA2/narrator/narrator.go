package narrator

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

func readCsvFile(filePath string) (records [][]string, err error) {

	f, err := os.Open(filePath)
	if err != nil {
		log.Println("Unable to read input file "+filePath, err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err = csvReader.ReadAll()

	if err != nil {
		log.Println("Unable to parse file as CSV for "+filePath, err)
	}

	return
}

type Narrator struct {
	Enabled bool

	characters map[int]string
}

func NewFromFile(path string) (n Narrator, err error) {

	n.Enabled = false
	n.characters = make(map[int]string)

	_, err = os.Stat(path)

	if errors.Is(err, os.ErrNotExist) {
		log.Println("Can't find narratorFile " + path)

		return
	}

	records, err := readCsvFile(path)

	if err != nil {

		return
	}

	for _, r := range records {

		// format: EPC,Name

		if len(r) < 2 {
			continue
		}

		id, err := strconv.Atoi(r[0])

		if err != nil {
			continue
		}

		n.characters[id] = r[1]
	}

	return
}

func (n *Narrator) SearchAndSay(id int) {

	character, ok := n.characters[id]

	if !ok {
		n.characters[id] = ""
		Say(strconv.Itoa(id))
	} else {
		Say(fmt.Sprintf("%s, Número: %d", character, id))
	}
}
