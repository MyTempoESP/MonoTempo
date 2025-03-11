#!/bin/sh

#git submodule update --remote --init --recursive

mkdir -p /usr/local/lib/

pip install -r requirements.txt

base=`pwd`
cd mytempo-chafon-reader

mkdir -p extern/chafon/lib
cp $base/default_lib/*.so ./extern/chafon/lib/
mkdir -p Release

cd Release

#cmake -DPRODUCTION=ON -DASSERTUTILS_RELEASE=ON -DPYTHON_BUILD=ON -DCMAKE_BUILD_TYPE=Release ..
cmake -DPRODUCTION=OFF -DASSERTUTILS_RELEASE=OFF -DPYTHON_BUILD=ON -DCMAKE_BUILD_TYPE=Debug ..

make && make install && cp python/pychafon.so ../..

cd ../..

#mkdir -p /usr/local/lib/assertutils
#mkdir -p /usr/local/lib/chafon_rfid
#mkdir -p /usr/local/lib/chafon

#find /usr/local/lib

#cp -r /usr/local/lib/assertutils lib/assertutils
#cp -r /usr/local/lib/chafon_rfid lib/chafon_rfid
#cp -r /usr/local/lib/chafon lib/chafon

#echo "export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib/assertutils:/usr/local/lib/chafon_rfid:/usr/local/lib/chafon" >> .venv/bin/activate

