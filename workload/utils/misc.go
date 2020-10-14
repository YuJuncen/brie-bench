package utils

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"strings"
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

type shortError struct{ error }

func (e shortError) String() string {
	return e.Error()
}

// ShortError prints a short error message.
func ShortError(err error) zap.Field {
	return zap.Stringer("error", shortError{err})
}
