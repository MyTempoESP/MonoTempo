package constant

import (
	"os"
	"time"
)

var (
	ProgramTimezone, _ = time.LoadLocation("Brazil/East")
	Version            = os.Getenv("PROGRAM_COMMIT_HASH")[:4]
	Reader             = os.Getenv("READER_NAME")
	Serie              = 501
	ReaderPath         = os.Getenv("READER_PATH")
	VersionNum         = os.Getenv("VERSION_NUMBER_AA2")
)

const (
	FORTH_BTN_PRESSED = "bac @ ba2 @ OR ."
)
