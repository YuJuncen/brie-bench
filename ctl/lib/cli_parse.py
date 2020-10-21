import re
from typing import Any, Callable, Iterable, List, NamedTuple, TypeVar
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
    ("cluster_version", str),
    ("component_version", dict)
])

version_flag = re.compile(r"--(tidb|tikv|pd)\.(version|hash)")

def shift() -> str:
    """
    shift shifts the current argv and return the first value. 
    Like `shift` command in bash.
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
        cluster_version = "nightly",
        component_version = {}
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
        elif tag == '--workload-storage':
            workload_storage = shift()
            conf.workload_storage = workload_storage
        elif tag == '--':
            conf.cargs = argv
            argv = []
        elif tag == '--cluster-version':
            conf.cluster_version = shift()
        elif re.match(version_flag, tag):
            matcher = re.match(version_flag, tag)
            component = matcher[1]
            spec_type = matcher[2]
            conf.component_version[component] = { spec_type: shift() }
        else:
            conf.other_args.append(tag)
    return conf

NOT_AN_NUMBER = "NAN"
USER_ABORT = "UAB"
def read_int():
    """
    read_int reads a integer from stdin,
    returns NOT_A_NUMBER if input cannot be parsed as int,
    returns USER_ABORT if user aborts.
    """
    try:
        return int(input())
    except EOFError:
        return USER_ABORT
    except KeyboardInterrupt:
        return USER_ABORT
    except:
        return NOT_AN_NUMBER

T = TypeVar('T')
def select(items: Iterable[T], mapper : Callable[[T], str] = lambda x: x) -> T:
    """
    select prompts user to select one item in the items.
    mapper maps the items to a printable string.
    """
    for i, item in enumerate(items):
        print(f"{i}) {mapper(item)}")
    print("#? ", end='')
    choice = read_int()
    while choice is NOT_AN_NUMBER or choice is USER_ABORT or choice >= len(items):
        if choice is USER_ABORT:
            print("(aborted)")
            exit(1) 
        print("(please, input a number in list)#? ", end='')
        choice = read_int()
    return items[choice]

def prompt_ok(prompt: str, default: bool = True) -> bool:
    """
    prompt_ok prompts the user for a boolean value.
    """
    hint = "[Y/n]" if default else "[y/N]"
    i = input(f"{prompt} {hint} ")
    if i.lower() == 'y':
        return True
    if i.lower() == 'n':
        return False
    if i == '':
        return default
    
    print("(please input y or n) ")
    return prompt_ok(prompt, default)