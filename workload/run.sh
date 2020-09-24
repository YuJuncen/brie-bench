#! /bin/bash

ARTIFICATS=/artifacts

set -eu

source ./utils.sh
debug_show_cluster_info

parse_args $@

component=${component-"none"}
workload=${workload-""}
hash=${hash-""}
components=(br)

case $component in
    "br" )
        ./br/build $hash
        ./br/run $workload
        ;;
    * )
        log "support components are ($components). Sorry for your choice $component is unsupported."
        ;;
esac