import logging
import requests
from typing import Any, List
from .cli_parse import Config, shift
from .saved_clusters import get_last_request, select_from_last_requests

def get_cluster(id: int, config: Config) -> Any:
    return requests.get(f"http://{config.api_server}/api/cluster/{id}").json()

def get_cluster_resources(id: int, config: Config) -> List[Any]:
    r = requests.get(f"http://{config.api_server}/api/cluster/resource/{id}")
    return r.json()

def get_metric_path(id: int, config: Config) -> str:
    resources = get_cluster_resources(id)
    for resource in resources:
        if "grafana" in resource["components"]:
            return resource["ip"]
    return ""

def run(config: Config):
    id = shift() or select_from_last_requests()
    if id == '.':
        id = get_last_request()

    get_type = shift() or "info"
    if get_type == "info":
        logging.info(get_cluster(int(id), config))
        logging.info(get_cluster_resources(int(id), config))
    if get_type == "metric":
        logging.info(get_metric_path(int(id), config))
    

