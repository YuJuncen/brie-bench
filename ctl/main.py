from typing import NoReturn
from lib.init_config import get_config, save_new_config
import coloredlogs
import logging
import sys
from lib.cli_parse import parse_cli
from lib.exec import ClusterRequest
from lib.get_cluster import run as get_cluster
from lib.get_file import FileReader

def exec_request():
    logging.basicConfig(format="%(asctime)s %(message)s", level=logging.INFO)
    config = get_config()
    cli_config = parse_cli(sys.argv)
    req = ClusterRequest(cli_config, config)
    req.request()

ops = {
    "exec": exec_request,
    "get_cluster": lambda: get_cluster(get_config()),
    "create_config": save_new_config,
    "get_file": FileReader(get_config()).run
}

def show_help(filename: str, exitcode:int=1) -> NoReturn :
    print(f"""usage: {filename} sub-command [...]
sub-commands:
    {", ".join(ops.keys())}""")
    exit(exitcode)


if __name__ == "__main__":
    coloredlogs.install()
    filename = sys.argv.pop(0)
    if len(sys.argv) == 0:
        show_help(filename)
    subcmd = sys.argv.pop(0)

    if ops.get(subcmd, None) is not None:
        ops.get(subcmd)()
    else:
        show_help(filename)
