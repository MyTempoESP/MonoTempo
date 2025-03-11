#!/bin/sh

pacman -S ntp

PATCHFILE=$(realpath ntp.conf.patch)

pushd .

cd /etc && patch < $PATCHFILE

popd

systemctl restart ntpd

timedatectl set-timezone America/Sao_Paulo
timedatectl set-ntp yes

echo "Done!"

