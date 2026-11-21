#!/bin/sh

sfdisk -W always /dev/sdb < MakeFAT32.txt && mkfs.vfat -F 32 -n MYTEMPO /dev/sdb1

