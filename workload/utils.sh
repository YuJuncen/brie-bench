debug_show_cluster_info() {
    output=${1-'/artifacts/cluster_info.txt'}
    echo "CLUSTER_ID: "   "${CLUSTER_ID-'UNKNOWN'}" >> $output
    echo "CLUSTER_NAME: " "${CLUSTER_NAME-'UNKNOWN'}" >> $output
    echo "TIDB_ADDR: "    "${TIDB_ADDR-'UNKNOWN'}" >> $output
    echo "PD_ADDR: "      "${PD_ADDR-'UNKNOWN'}" >> $output
    echo "PROM_ADDR: "    "${PROM_ADDR-'UNKNOWN'}" >> $output
    echo "API_SERVER: "   "${API_SERVER-'UNKNOWN'}" >> $output
}

export BLUE_FONT=$'\e[36m'
export RED_FONT=$'\e[31m'
export RESET_FONT=$'\e[0m'

fail() {
    echo -e "$RED_FONT[$(date "+%Y-%m-%d %H:%M:%S")]$RESET_FONT" $@
    exit 1
}

log() {
    if [ ${logfile-""} ]; then
        echo "[$(date "+%Y-%m-%d %H:%M:%S")]" $@ >> $logfile
    fi

    echo -e "$BLUE_FONT[$(date "+%Y-%m-%d %H:%M:%S")]$RESET_FONT" $@
}

show_help() {
    cat <<eof
usage:
    exec component [-h/--hash (BR commit hash)] [-w workload]

    library needed:
        workload/utils.sh
eof
}

parse_args() {
    export component=${1-""}
    shift
    while [[ $# -gt 0 ]]; do
    case $1 in
    -h | --hash )
        export hash=$2
        shift
        shift
        ;;
    -w )
        export workload=$2
        shift
        shift
        ;;
    *)
        show_help
        fail "unknown flag $1"
        ;;
    esac
    done
}