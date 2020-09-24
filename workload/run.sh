#! /bin/bash

ARTIFICATS=/artifacts

set -eu

source ./utils.sh
debug_show_cluster_info

target=${target-"none"}
workload=${workload-""}
hash=${hash-""}
components=(br)

case $target in
    "br" )
        ./br/build $hash
        ./br/run $workload
        ;;
    * )
        log "please set the target env to one of ($components). Your choice is $target."
        ;;
esac