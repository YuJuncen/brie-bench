all: workload ctl

prerequest:
	@which python3 || (echo "python3 binary not found, please install python firstly." && exit 1)

docker-bulid:
	docker build --tag lovecsust/brie-bench .

workload-image:
	go build -o bin/run workload/main.go

venv-prepare:
	python3 -m venv ctl/venv

ctl-build: prerequest venv-prepare
	$(VENV_ENABLE)
	pip install -r ctl/requirements.txt
	make generate_config

generate_config:
	python3 ctl/main.py create_config
	