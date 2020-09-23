debug_show_cluster_info() {
    output=${1-'/artifacts/cluster_info.txt'}
    echo "CLUSTER_ID: " "$CLUSTER_ID" >> $output
    echo "CLUSTER_NAME: " "$CLUSTER_NAME" >> $output
    echo "TIDB_ADDR: " "$TIDB_ADDR" >> $output
    echo "PD_ADDR: " "$PD_ADDR" >> $output
    echo "PROM_ADDR: " "$PROM_ADDR" >> $output
    echo "API_SERVER: " "$API_SERVER" >> $output
}