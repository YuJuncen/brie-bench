#! /bin/bash

git clone https://github.com/pingcap/br && cd br
make build
if [ $hash ]; then
    git reset --hard $hash
fi

bin/br -V > /artifacts/output.txt