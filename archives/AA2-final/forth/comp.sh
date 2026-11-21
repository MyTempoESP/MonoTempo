#!/bin/sh

rm -i ../Arduino/aa2/forth2cstr.h
rm -i ../Arduino/aa2/.forth2cstr.h.swp

which vim 1> /dev/null 2> /dev/null

if [ $? -eq 1 ] ; then
	echo "This script requires vim."
	exit 1
fi

vim --clean -c "source unforth.vim" -c "edit aa2.fs" -c "Unforth ../Arduino/aa2/forth2cstr.h" -c "q" 1> /dev/null 2> /dev/null

