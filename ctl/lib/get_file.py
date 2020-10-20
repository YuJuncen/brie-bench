import logging
from .cli_parse import Config, select, shift
from .saved_clusters import from_cli_desc
from typing import List
from colorama import Fore
from minio import Minio
from minio.definitions import Object
import subprocess
import sys


# TODO allow artifacts bucket configurable
artifacts_bucket = "artifacts"

def last_part(path: str) -> str :
    """
    last_part returns the last part of the path.

    example:
    /foo/bar/baz/ -> baz/
    /foo/bar -> bar
    """
    if path.endswith('/'):
        return path.split('/')[-2] + '/'
    return path.split('/')[-1]

class FileReader:
    def __init__(self, config: Config):
        self._config = config
        self._client = Minio(config.s3_endpoint,
            access_key=config.s3_access_key,
            secret_key=config.s3_secret_key,
            secure=False)
    
    def get_artifacts_dir_of(self, cluster_id: int) -> Object:
        objs = list(self._client.list_objects(artifacts_bucket, f"{cluster_id}/"))
        if len(objs) == 0:
            return None
        if len(objs) > 1:
            return select(objs, lambda o: f"{o.object_name.strip('/')}")
        return objs[0]

    def select_file(self, dir: str) -> Object:
        dir = f"{dir}/" if not dir.endswith("/") else dir
        objs : List[Object] = list(self._client.list_objects(artifacts_bucket, dir))
        print(f"current path: {Fore.GREEN}{dir}{Fore.RESET}")
        new_obj = select(objs, lambda o: f"{last_part(o.object_name)}")
        if new_obj.is_dir:
            return self.select_file(f"{new_obj.object_name}")
        return new_obj
        
    def query_file(self, cluster_id: int):
        artifact = self.get_artifacts_dir_of(cluster_id)
        if artifact is None:
            logging.error(f"the cluster {cluster_id} seems has no artifacts, maybe it's doesn't end now?")
            exit(1)
        obj = self.select_file(artifact.object_name)
        data = self._client.get_object(obj.bucket_name, obj.object_name)
        proc = subprocess.Popen(["less"], stdin = subprocess.PIPE)
        # this would load the total file into memory, if the file is big,
        # we might meet some error, but don't worry for now.
        proc.communicate(data.read())

    def run(self):
        cluster = from_cli_desc(shift())
        file = shift()
        if file is None:
            self.query_file(cluster)
            return
        artifact = self.get_artifacts_dir_of(cluster)
        obj = self._client.get_object(artifacts_bucket, f"{artifact.object_name}{file}")
        for d in obj.stream(32 * 1024):
            sys.stdout.write(d.decode('utf-8'))
        



    