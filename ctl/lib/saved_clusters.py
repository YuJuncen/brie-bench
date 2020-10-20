import logging
import tempfile
from .cli_parse import select
from typing import List
from pathlib import Path

saved_file = ".brie-bench-requested-clusters"
saved_path = Path.home() / saved_file

def all_requests() -> List[str]:
    with open(saved_path, "r") as f:
        return [line for line in f]

def get_id_of_record(record: str) -> int:
    id_idx = record.index(':')
    return int(record[:id_idx])

def save_request_id(id: int, cmd: str):
    with open(saved_path, "a") as f:
        f.write(f"{id}:{cmd}")

def get_last_request() -> int:
    with open(saved_path, "r") as f:
        *_, last_line = f
        return get_id_of_record(last_line)

def select_from_last_requests() -> int:
    requests = all_requests()
    return get_id_of_record(select(requests))

def from_cli_desc(cli: str) -> int:
    """
    from_cli_desc parse from cli args.

    None or empty string: prompt the user for one
    a dot(.): get the last one
    number: the cluster ID
    """
    if cli is None or cli == "":
        return select_from_last_requests()
    if cli == ".":
        return get_last_request()
    try:
        return int(cli)
    except:
        logging.error(f"failed to parse your input! please input one of '.' or number. (you input {cli})")
        exit(1)
    
    
