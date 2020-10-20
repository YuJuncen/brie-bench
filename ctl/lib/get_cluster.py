import logging
import sys
import requests
import json
from typing import Any, List
from .cli_parse import Config, shift
from .saved_clusters import from_cli_desc

def get_cluster(id: int, config: Config) -> Any:
    return requests.get(f"http://{config.api_server}/api/cluster/{id}").json()

def get_cluster_resources(id: int, config: Config) -> List[Any]:
    r = requests.get(f"http://{config.api_server}/api/cluster/resource/{id}")
    return r.json() or []

def get_metric_path(id: int, config: Config) -> str:
    resources = get_cluster_resources(id, config)
    for resource in resources:
        if "grafana" in resource["components"]:
            return resource["ip"]
    return ""

def run(config: Config):
    id = from_cli_desc(shift())
    get_type = shift() or "info"
    if get_type == "info":
        json.dump(get_cluster(int(id), config), sys.stdout)
    elif get_type == "resource":
        json.dump(get_cluster_resources(int(id), config), sys.stdout)
    elif get_type == "metric" or get_type == "grafana":
        addr = get_metric_path(id, config)
        if addr == "":
            logging.error("cannot find cluster metrics")
            exit(1)
        print(f"http://{get_metric_path(id, config)}:3000")
    else:
        logging.error(f"unsupported cluster resource type '{get_type}'")
    

