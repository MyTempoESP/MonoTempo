package constant

import (
	"os"
	"time"
)

var (
	ProgramTimezone, _ = time.LoadLocation("Brazil/East")
	Reader             = os.Getenv("READER_NAME")
	DeviceId           = os.Getenv("MYTEMPO_DEVID")
	ReaderPath         = os.Getenv("READER_PATH")
	VersionNum         = os.Getenv("VERSION_NUMBER_AA2")
	SerialPortOverride = os.Getenv("SERIAL_PORT_OVERRIDE")
	ReaderType         = os.Getenv("MYTEMPO_READER_TYPE")
)
