package lcdlogger

import (
	"strconv"
	"strings"
)

type ForthNumber struct {
	Value     int64
	Magnitude int // 1, 10, 100, 1000 (10^Magnitude)
}

func IPIfy(ip string) (out IPOctets) {

	out = IPOctets{0, 0, 0, 0}

	parts := strings.Split(ip, ".")

	if len(parts) != 4 {

		return
	}

	var err error
	var num int

	for i, part := range parts {

		num, err = strconv.Atoi(part)

		if err != nil {

			return
		}

		if num < 0 || num > 255 {

			return
		}

		out[i] = num
	}

	return
}

func ToForthNumber(n int64) (f ForthNumber) {

	if n < 1000 {

		f.Value = n
		f.Magnitude = 0

		return
	}

	if n < 1_000_000 {

		f.Value = n / 1000
		f.Magnitude = 3

		return
	}

	f.Value = n / 1_000_000
	f.Magnitude = 6

	return
}
