all: workload ctl

CLEAN := echo clean

docker-bulid:
	docker build --tag lovecsust/brie-bench .

workload-image:
	go build -o bin/run workload/main.go

ctl:
	echo "ctl isn't ready yet, exiting :(" && exit 1