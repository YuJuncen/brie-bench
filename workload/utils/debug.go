package utils

import (
	"fmt"
	"github.com/yujuncen/brie-bench/workload/config"
	"io"
	"os"
	"path"
)

const (
	ClusterInfoFile = "cluster_info.txt"
	EnvInfoFile     = "env.txt"
	infoFormat      = `Cluster ID: %20s
Cluster name: %20s
PD address: %20s
TiDB address: %20s
Prometheus address: %20s
API server address: %20s
`
)

func DumpCluster(c *BenchContext) error {
	file, err := os.Create(path.Join(config.Artifacts, ClusterInfoFile))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, infoFormat, c.id, c.name, c.PdAddr, c.TidbAddr, c.PrometheusAddr, c.apiAddr)
	return err
}

func DumpEnvTo(file io.Writer) error {
	for _, env := range os.Environ() {
		if _, e := fmt.Fprintln(file, env); e != nil {
			return e
		}
	}
	return nil
}

func DumpEnv() error {
	file, err := os.Create(path.Join(config.Artifacts, EnvInfoFile))
	if err != nil {
		return err
	}
	return DumpEnvTo(file)
}
