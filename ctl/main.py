import coloredlogs
import logging
import sys
from lib.cli_parse import default_config, parse_cli
from lib.exec import ClusterRequest
from lib.get_cluster import run as get_cluster

def exec_request():
    logging.basicConfig(format="%(asctime)s %(message)s", level=logging.INFO)
    config = default_config()
    cli_config = parse_cli(sys.argv)
    req = ClusterRequest(cli_config, config)
    req.request()

def show_help(filename: str, exitcode:int=1):
    print(f"""usage: {filename} subcmd [...]
subcmds:
    exec
    get_file
    get_cluster
    rebuild_metrics
""")
    exit(exitcode)

if __name__ == "__main__":
    coloredlogs.install()
    filename = sys.argv.pop(0)
    if len(sys.argv) == 0:
        show_help(filename)
    subcmd = sys.argv.pop(0)
    if subcmd == "exec":
        exec_request()
    elif subcmd == "get_cluster":
        config = default_config()
        get_cluster(config)
    else:
        show_help(filename)
