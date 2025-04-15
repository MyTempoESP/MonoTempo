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
	WifiIface          = os.Getenv("WIFI_IFACE")
	Lte4GIface         = os.Getenv("LTE4G_IFACE")
)
