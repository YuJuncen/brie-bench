FROM golang:1.14-alpine as builder
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct

WORKDIR /brie
COPY ./workload .

RUN go build -o bin/run

FROM golang:1.14-alpine
WORKDIR /brie
ENV ARTIFICATS=/artifacts
ENV hash=""
ENV workload=""
ENV target="none"
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct

RUN apk add --no-cache \
    make \
    git \
    bash \
    curl \
    gcc \
    g++
RUN mkdir $ARTIFICATS
COPY --from=builder /brie/bin/run bin/run

CMD ["bin/run", "--help"]

