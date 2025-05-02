package main

import (
	"encoding/csv"
	"os"

	"go.uber.org/zap"
)

func CreateCSVWriter(filename string) (*csv.Writer, *os.File, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, nil, err
	}
	writer := csv.NewWriter(f)
	return writer, f, nil
}

func WriteCSVRecord(writer *csv.Writer, record []string, logger *zap.Logger) {
	err := writer.Write(record)
	if err != nil {
		logger.Warn("error writing csv record", zap.Error(err))
	}
}
