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
export RESET_FONT=$'\e[0m'

log() {
    if [ ${logfile-""} ]; then
        echo "[$(date "+%Y-%m-%d %H:%M:%S")]" $@ >> $logfile
    fi

    echo -e "$BLUE_FONT[$(date "+%Y-%m-%d %H:%M:%S")]$RESET_FONT" $@
}