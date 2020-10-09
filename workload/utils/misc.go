package utils

import (
	"errors"
	"fmt"
	"github.com/pingcap/log"
	"go.uber.org/zap"
	"strings"
	"time"
)

var NOPIO = NIO{}

type NIO struct{}

func (n2 NIO) Read(p []byte) (n int, err error) {
	return len(p), nil
}

func (n2 NIO) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func HostAndPort(addr string) (host string, port string, err error) {
	delim := strings.LastIndex(addr, ":")
	if delim == -1 {
		err = errors.New(fmt.Sprintf("bad input host and port format %v", addr))
		return
	}
	host, port = addr[:delim], addr[delim+1:]
	return
}

// Bench runs the task, with logging the time cost.
func Bench(name string, task func() error) error {
	start := time.Now()
	defer func() {
		log.Info("bench task done", zap.String("name", name), zap.Duration("cost", time.Since(start)))
	}()
	return task()
}
