import logging
import requests
from typing import Any, List
from .cli_parse import Config, shift
from .saved_clusters import from_cli_desc

def get_cluster(id: int, config: Config) -> Any:
    return requests.get(f"http://{config.api_server}/api/cluster/{id}").json()

def get_cluster_resources(id: int, config: Config) -> List[Any]:
    r = requests.get(f"http://{config.api_server}/api/cluster/resource/{id}")
    return r.json()

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
        logging.info(get_cluster(int(id), config))
        logging.info(get_cluster_resources(int(id), config))
    if get_type == "metric":
        logging.info(f"view grafana at http://{get_metric_path(int(id), config)}:3000")
    

