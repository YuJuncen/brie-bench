#! /bin/bash

ARTIFICATS=/artifacts

source ./utils.sh
debug_show_cluster_info

target=${target-"none"}
workload=${workload-""}
hash=${hash-""}

case $target
    "br" )
        ./br/build 
        ./br/run $workload
        ;;
    "none" )
        echo "please set the target env to one of (br)."