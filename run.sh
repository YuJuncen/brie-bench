#! /bin/bash

ARTIFICATS=/artifacts

source ./utils.sh
debug_show_cluster_info

git clone https://github.com/pingcap/br && cd br
if [ $hash ]; then
    git reset --hard $hash
fi

make build

bin/br restore full \
    --log-file ${ARTIFICATS-.}/br-restore-full.log \
    --log-level DEBUG \
    --pd \
    -s 's3://mybucket/ycsb?access-key=YOURACCESSKEY&secret-access-key=YOURSECRETKEY' \
    --s3.endpoint 'http://172.16.4.4:30812'