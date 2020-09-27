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
        ./utils.sh
eof
}

parse_args() {
    export component=${1-""}
    shift
    export other_args="[]"
    while [[ $# -gt 0 ]]; do
    case $1 in
    --dry-run )
        export dry_run=1
        shift
        ;;
    *)
        other_args=`echo $other_args | jq ". += [\"$1\"]"`
        shift
        ;;
    esac
    done
}

add_cluster() {
  echo $1 >> .brie_bench_last_cluster
}

get_cluster() {
  clusters=$(cat .brie_bench_last_cluster || echo "")
  if [ ! "$clusters" ]; then
    fail "no request found"
  fi
  select cluster in $clusters; do
    echo $cluster
    break
  done
}