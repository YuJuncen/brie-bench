#! /bin/bash

ARTIFICATS=/artifacts

set -eu

source ./utils.sh
debug_show_cluster_info

parse_args $@

target=${component-"none"}
workload=${workload-""}
hash=${hash-""}
components=(br)

case $target in
    "br" )
        ./br/build $hash
        ./br/run $workload
        ;;
    * )
        log "support components are ($components). Sorry for your choice $target is unsupported."
        ;;
esac