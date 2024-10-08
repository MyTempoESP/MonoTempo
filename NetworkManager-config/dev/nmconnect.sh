#!/bin/sh

while true; do
	BUF=$(cat /etc/wifi-connection)

	printf "\033[32;1mAttempting to connect: $NETWORK\033[0m"

	sed -i "/^ssid=/ s/=.*\$/=$NETWORK/" /etc/NetworkManager/system-connections/Wifi.nmconnection
	nmcli reload
done
