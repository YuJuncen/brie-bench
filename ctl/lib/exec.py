from .saved_clusters import save_request_id
import requests
import json
import logging

from requests.sessions import Request
from .cli_parse import CLIConfig, Config, default_config, parse_cli, s3_args 


def tikv(node_id: int) -> dict:
    return {"component": "tikv", "deploy_path": "/data1", "rri_item_id": node_id}

class ClusterRequest:
    def __init__(self, cli_config: CLIConfig, config: Config):
        self.cli_config = cli_config
        self.config = config
        logging.info(f"cli config {self.cli_config}")
        self.data = {
            "cluster_request": {
                "name": "brie-test",
                "version": "nightly"
            },
            "cluster_request_topologies": [
                {"component": "tidb", "deploy_path": "/data1", "rri_item_id": 1}, 
                {"component": "pd", "deploy_path": "/data1", "rri_item_id": 1}, 
                {"component": "prometheus", "deploy_path": "/data1", "rri_item_id": 1},
                {"component": "grafana", "deploy_path": "/data1", "rri_item_id": 1}
            ],
            "cluster_workload": {        
                "type": "standard",
                "docker_image": "lovecsust/brie-bench:latest",
                "cmd": "bin/run",
                "args": [],
                "rri_item_id": 1,                      
                "artifact_dir": "/artifacts",
                "env": {}
            }
        }
        self.set_component(self.cli_config.component)
        if cli_config.pr_mode:
            self.pr_mode()
        if cli_config.import_to != "":
            self.import_mode(cli_config.import_to)
        self.set_workload(self.cli_config.workload, self.cli_config.workload_storage)
        self.set_tikv_count(3)
        self.init_args()
    
    def init_args(self):
        """
        init_args initializes arguments for the component and the workload image.
        """
        self.__add_args(*self.cli_config.other_args)
        self.component_args(*self.cli_config.cargs)
        self.data["cluster_request"]["name"] = f"{self.cli_config.component}-bench-{self.cli_config.workload}"
        return self

    def set_workload(self, workload_name: str, workload_storage: str = ""):
        self.workload = workload_name
        self.__add_args("--workload-name", workload_name)
        workload_storage = workload_storage or f"s3://{self.config.workload_bucket}/{self.cli_config.workload}?{s3_args(self.config)}"
        self.__add_args("--workload-storage", workload_storage)        

    def set_tikv_count(self, count : int = 1):
        for i in range(2, 2 + count):
            self.data["cluster_request_topologies"].append(tikv(i))
        return self
    
    def pr_mode(self):
        self.data["cluster_workload"]["type"] = "pr"
        return self

    def __add_args(self, *args: str):
        for arg in args:
            self.data["cluster_workload"]["args"].append(arg)

    def import_mode(self, to: str):
        self.data["cluster_workload"]["type"] = "importer"
        self.data["cluster_workload"]["backup_path"] = f"{self.config.workload_bucket}/{to}"
        return self

    def component_args(self, *args: str):
        for arg in args:
            self.__add_args("--cargs", arg)
        
    def set_component(self, component: str):
        self.__add_args("--component", component)
        if component == "dumpling":
            self.data["cluster_workload"]["restore_path"] = f"${self.config.workload_bucket}/${self.cli_config.workload}"
        return self

    def json(self) -> str:
        return json.dumps(self.data)

    def request(self):
        try:
            logging.info(f"request with json {self.json()}")
            if not self.cli_config.dry_run:
                resp = requests.post(f"http://{self.config.api_server}/api/cluster/test", json=self.data)
                if resp.status_code / 100 != 2:
                    logging.error(f"failed to request: got response with failed message {resp}: {resp.content}")
                resp_json = resp.json()
                cluster_id = resp_json['cluster_request_id']
                logging.info(f"success request, cluster ID {cluster_id}")
                save_request_id(cluster_id, f"{self.cli_config.component} {' '.join(self.cli_config.cargs)}")
        except Exception as e:
            logging.error(f"failed to request: {e}")


