from typing import Iterable, List, NamedTuple
from recordclass import recordclass
import logging
import sys

Config = NamedTuple("Config", [
    ("s3_endpoint", str),
    ("s3_access_key", str),
    ("s3_secret_key", str),
    ("workload_bucket", str),
    ("api_server", str),
])

# TODO remove this function and query for config.
def default_config() -> Config:
    return Config(
        s3_endpoint = "172.16.4.4:30812",
        s3_access_key = "YOURACCESSKEY",
        s3_secret_key = "YOURSECRETKEY",
        workload_bucket = "brie",
        api_server = "172.16.5.110:8000"
    )

def s3_args(conf: Config) -> str:
    """
    s3_args make the query string for BRIE style storage URI.
    """
    # BRIE requires the http[s]:// prefix; currently, we don't support TLS, just put http:// directly.
    # force-path-style is needed since endpoint probably be a IP address for testing environment.
    return f"access-key={conf.s3_access_key}&secret-access-key={conf.s3_secret_key}&endpoint=http://{conf.s3_endpoint}&force-path-style=true"

CLIConfig = recordclass("CLIConfig", [
    ("other_args", list), 
    ("cargs", list),
    ("import_to", str),
    ("pr_mode", bool),
    ("component", str),
    ("dry_run", bool),
    ("workload", str),
    ("workload_storage", str),
])

def shift() -> str:
    """
    shift shifts the current argv. Like `shift` command in bash.
    """
    try:
        return sys.argv.pop(0)
    except:
        return None

def parse_cli(argv: List[str], conf : CLIConfig = None) -> CLIConfig:
    """
    parse_cli parses the input CLI args.
    """
    conf = conf or CLIConfig(
        other_args = [],
        cargs = [],
        import_to = "",
        pr_mode = False,
        component = "",
        dry_run = False,
        workload = "",
        workload_storage = "",
    )
    if len(argv) == 0:
        logging.error("please specify the component name")
        exit(1)

    conf.component = shift()

    while len(argv) > 0:
        tag = shift()
        if tag == '--import-to':
            conf.import_to = shift()
        elif tag == '--pr':
            conf.pr_mode = True
        elif tag == '--dry-run':
            conf.dry_run = True
        elif tag == '--workload-name':
            workload_name = shift()
            conf.workload = workload_name
            conf.other_args.extend(["--workload-name", workload_name])
        elif tag == '--workload-storage':
            workload_storage = shift()
            conf.workload_storage = workload_storage
            conf.other_args.extend(["--workload-storage", workload_storage])
        elif tag == '--':
            conf.cargs = argv
            argv = []
        else:
            conf.other_args.append(tag)
    return conf
        
def read_int():
    try:
        return int(input())
    except:
        return None

def select(items: Iterable[str]) -> str:
    for i, item in enumerate(items):
        print(f"{i}) {item}")
    print("#? ")
    choice = read_int()
    while choice is None or choice >= len(items):
        print("(please, input a number in list)#? ")
        choice = read_int()
    return items[choice]