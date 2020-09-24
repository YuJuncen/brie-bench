FROM golang:1.14-alpine
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct
ENV ARTIFICATS=/artifacts
ENV hash=""
ENV workload=""
ENV target="none"

WORKDIR /brie
COPY ./workload .

RUN apk add --no-cache \
    make \
    git \
    bash \
    curl \
    gcc

RUN mkdir $ARTIFICATS
CMD ["./run.sh"]