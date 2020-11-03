from minio.api import Minio
from .cli_parse import Config, default_config, prompt_ok, select
from getpass import getpass
from pathlib import Path
from colorama import Fore
import logging
import json

config_path = Path.home() / '.config'
config_file = Path.home() / '.config' / 'brie-bench-config.json'

def show_config(conf : Config):
    for key, value in conf._asdict().items():
        if key == 's3_secret_key':
            print(f"{Fore.BLUE}[ %-16s ]{Fore.RESET} %s" % (key, '*' * len(value)))
        else:
            print(f"{Fore.BLUE}[ %-16s ]{Fore.RESET} %s" % (key, value))

def prompt_config() -> Config:
    print("current config: ")
    old_config = get_config()
    show_config(old_config)
    print(f"{Fore.GREEN}you can leave any field empty to leave this field untouched then.{Fore.RESET}")
    api_server = (input("Please input your api server: ") or old_config.api_server).strip()
    s3_endpoint = (input("Please input your s3 endpoint (without http:// prefix, TLS isn't supported for now): ") or old_config.s3_endpoint).strip()
    s3_access_key = input("Please input your s3 access key: ") or old_config.s3_access_key
    s3_secret_key = getpass("Please input your s3 secret access key: ") or old_config.s3_secret_key
    logging.info("dialing to s3 storage using the current config...")
    try:
        client = Minio(s3_endpoint, s3_access_key, s3_secret_key, secure=False)
        names = map(lambda b: b.name, client.list_buckets())
        print(f"Which bucket you store your workloads?")
        # leave empty string here for not changed. mapping it to a user-friend string
        workload_bucket = select(["", *names], lambda b: b or "<not changed>") or old_config.workload_bucket
    except Exception as e:
        cont = prompt_ok(f"{Fore.RED}cannot connect to s3 ({e}){Fore.RESET}, continue?", False)
        if not cont:
            exit(1)
        workload_bucket = input("What bucket you store your workloads? ") or old_config.workload_bucket
    
    return Config(
        api_server = api_server,
        s3_endpoint = s3_endpoint,
        s3_access_key = s3_access_key,
        s3_secret_key = s3_secret_key,
        workload_bucket = workload_bucket
    )

def save_new_config():
    conf = prompt_config()
    config_path.mkdir(parents=True, exist_ok=True)
    with open(config_file, "w") as f:
        json.dump(conf._asdict(), f)
    print(f"{Fore.GREEN}config saved!{Fore.RESET} (@ {config_file})")
    show_config(conf)

def get_config():    
    try:
        with open(config_file, "r") as f:
            return Config(**json.load(f))
    except Exception as e:
        logging.warn(f"exception({e}) during read config from file, using deafult config")
        return default_config()