#! /bin/bash

ARTIFICATS=/artifacts

set -eu

source ./utils.sh
debug_show_cluster_info

target=${target-"none"}
workload=${workload-""}
hash=${hash-""}

case $target in
    "br" )
        ./br/build $hash
        ./br/run $workload
        ;;
    "none" )
        echo "please set the target env to one of (br)."
        ;;
esac