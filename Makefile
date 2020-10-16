all: workload ctl


docker-bulid:
	docker build --tag lovecsust/brie-bench .

workload-image:
	go build -o bin/run workload/main.go

venv-prepare:
	python3 -m venv ctl/venv

ctl-build: venv-prepare
	$(VENV_ENABLE)
	pip install -r ctl/requirements.txt
	