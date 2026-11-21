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
)

const (
	FORTH_BTN_PRESSED = "btn ."
)
