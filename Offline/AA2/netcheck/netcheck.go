package netcheck

import (
	"fmt"
	"net"
)

// Check if given interface is up and has a meaningful ip address
func CheckIface(ifaceName string) (bool, error) {

	iface, err := net.InterfaceByName(ifaceName)

	if err != nil {
		return false, fmt.Errorf("interface not found: %w", err)
	}

	if iface.Flags&net.FlagUp == 0 {

		return false, nil
	}
	if iface.Flags&net.FlagLoopback != 0 {
		return false, nil
	}

	addrs, err := iface.Addrs()

	if err != nil {
		return false, fmt.Errorf("failed to get interface addresses: %w", err)
	}

	for _, addr := range addrs {

		if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.IsGlobalUnicast() {
			return true, nil
		}
	}

	return false, nil
}
