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

var iecUnits = []string{"", "K", "M", "G", "T", "P"}

// ToIce convert a size to iec format("xxxM")
func ToIec(num int64) string {
	unitIdx := 0
	floatNum := float64(num)
	for floatNum > 1024 && unitIdx < len(iecUnits)-1 {
		floatNum = floatNum / 1024
		unitIdx += 1
	}
	return fmt.Sprintf("%.1f%s", floatNum, iecUnits[unitIdx])
}

type Size uint64

func (s Size) String() string {
	return ToIec(int64(s))
}

func DumpCluster(c *Cluster) error {
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
