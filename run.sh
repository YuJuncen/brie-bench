#! /bin/bash

git clone https://github.com/pingcap/br && cd br
make build

bin/br -V > /artifacts/output.txt