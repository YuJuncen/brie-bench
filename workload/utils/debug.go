package utils

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/yujuncen/brie-bench/workload/config"
)

const (
	// ClusterInfoFile is the file for storing cluster info.
	ClusterInfoFile = "cluster_info.txt"
	// EnvInfoFile is the file for storing environment variables.
	EnvInfoFile = "env.txt"
	// infoFormat is the format for cluster info.
	infoFormat = `Cluster ID: %20s
Cluster name: %20s
PD address: %20s
TiDB address: %20s
Prometheus address: %20s
API server address: %20s
`
)

// DumpCluster dumps the cluster by the context to ClusterInfoFile.
func DumpCluster(c *BenchContext) error {
	file, err := os.Create(path.Join(config.Artifacts, ClusterInfoFile))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, infoFormat, c.id, c.name, c.PdAddr, c.TidbAddr, c.PrometheusAddr, c.apiAddr)
	return err
}

// DumpEnvTo dump environment vairbales to the writer.
func DumpEnvTo(file io.Writer) error {
	for _, env := range os.Environ() {
		if _, e := fmt.Fprintln(file, env); e != nil {
			return e
		}
	}
	return nil
}

// DumpEnv dumps the environment vairbales to the EnvInfoFile.
func DumpEnv() error {
	file, err := os.Create(path.Join(config.Artifacts, EnvInfoFile))
	if err != nil {
		return err
	}
	return DumpEnvTo(file)
}
