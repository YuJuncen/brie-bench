package utils

import (
	"fmt"
	"io"
	"strings"

	"go.uber.org/zap"
)

// HostAndPort parses the host and part for a socket address string.
func HostAndPort(addr string) (host string, port string, err error) {
	delim := strings.LastIndex(addr, ":")
	if delim == -1 {
		err = fmt.Errorf("bad input host and port format %v", addr)
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

// CombineWriters creates a combined io.Writer like io.MultiWriter,
// but skips nil parameters.
func CombineWriters(ws ...io.Writer) io.Writer {
	writers := make([]io.Writer, 0, len(ws))
	for _, w := range ws {
		if w != nil {
			writers = append(writers, w)
		}
	}
	return io.MultiWriter(writers...)
}
