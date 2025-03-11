package pinger

import (
	"errors"
	"net"
)

func Get4GIP() (string, error) {

	iface, err := net.InterfaceByName("enp0s21f0u3u3")

	if err != nil {

		return "", err
	}

	if iface.Flags&net.FlagUp == 0 {

		return "", errors.New("interface is down")
	}

	addrs, err := iface.Addrs()

	if err != nil {

		return "", err
	}

	for _, addr := range addrs {

		var ip net.IP

		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if ip == nil || ip.IsLoopback() {

			continue
		}

		ip = ip.To4()

		if ip == nil {

			continue // not an ipv4 address
		}

		return ip.String(), nil
	}

	return "", nil
}
