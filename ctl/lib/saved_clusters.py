import tempfile
import os.path
from .cli_parse import select
from typing import List

saved_file = ".brie-bench-requested-clusters"
saved_path = os.path.join(tempfile.gettempdir(), saved_file)

def all_requests() -> List[str]:
    with open(saved_path, "a") as f:
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
    
    
