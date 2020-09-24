debug_show_cluster_info() {
    output=${1-'/artifacts/cluster_info.txt'}
    echo "CLUSTER_ID: "   "${CLUSTER_ID-'UNKNOWN'}" >> $output
    echo "CLUSTER_NAME: " "${CLUSTER_NAME-'UNKNOWN'}" >> $output
    echo "TIDB_ADDR: "    "${TIDB_ADDR-'UNKNOWN'}" >> $output
    echo "PD_ADDR: "      "${PD_ADDR-'UNKNOWN'}" >> $output
    echo "PROM_ADDR: "    "${PROM_ADDR-'UNKNOWN'}" >> $output
    echo "API_SERVER: "   "${API_SERVER-'UNKNOWN'}" >> $output
}

log() {
    echo -e "\e[36m[$(date "+%Y-%m-%d %H:%M:%S")]\e[0m" $@
}