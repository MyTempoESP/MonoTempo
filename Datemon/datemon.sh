#!/bin/sh
#
# XXX: USE OF LOCALTIME
while true; do

	echo "Pinging $READER_IP"

	if ping $READER_IP; then
		echo "Setting reader clock: $(date +%H:%M:%S)"
		(sleep 1; echo root; sleep 1; echo impinj; sleep 1; echo "config system time $(date +'%Y.%m.%d-%H:%M:%S')"; sleep 1; echo exit) | telnet $READER_IP
	fi

	sleep 40 # wow
done
