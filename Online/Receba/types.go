package main

import (
	"database/sql"
)

type Receba struct {
	db *sql.DB

	AtletasRota string
	DeviceRota  string
	ProvaRota   string
	StaffRota   string
	InfoRota    string
}
