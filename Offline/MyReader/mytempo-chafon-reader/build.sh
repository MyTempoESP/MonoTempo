#!/bin/sh

cd Release

cmake -DPRODUCTION=ON -DASSERTUTILS_RELEASE=ON -DPYTHON_BUILD=ON -DCMAKE_BUILD_TYPE=Release ..
make && sudo make install

