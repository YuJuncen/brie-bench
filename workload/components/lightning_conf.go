package component

import (
	"fmt"
	"github.com/pingcap/log"
	"github.com/yujuncen/brie-bench/workload/config"
	"go.uber.org/zap"
	"io"
	"os"
	"path"
)

const (
	basic = `[lightning]
index-concurrency = %d
table-concurrency = %d
IO-concurrency = %d
`
	local = `[tikv-importer]
backend = "local"
region-split-size = %d
send-kv-pairs = %d
sorted-kv-dir = "%s"
range-concurrency = %d
`
)

type LightningConf struct {
	IndexConcurrency uint
	TableConcurrency uint
	IOConcurrency    uint
}

func NewLightningConf() *LightningConf {
	return &LightningConf{
		IndexConcurrency: 2,
		TableConcurrency: 6,
		IOConcurrency:    5,
	}
}

func (l *LightningConf) To(w io.Writer) error {
	_, err := fmt.Fprintf(w, basic, l.IndexConcurrency, l.TableConcurrency, l.IOConcurrency)
	return err
}

type LocalBackendConf struct {
	RegionSplitSize  uint
	SendKVPairs      uint
	SortedKVDir      string
	RangeConcurrency uint
}

func NewLocalBackendConf(dir string) *LocalBackendConf {
	return &LocalBackendConf{
		RegionSplitSize:  1024 * 1024 * 96,
		SendKVPairs:      32768,
		SortedKVDir:      dir,
		RangeConcurrency: 16,
	}
}

func (conf *LocalBackendConf) To(w io.Writer) error {
	_, err := fmt.Fprintf(w, local, conf.RegionSplitSize, conf.SendKVPairs, conf.SortedKVDir, conf.RangeConcurrency)
	return err
}

type LightningConfigFile struct {
	Path string
	IO   io.WriteCloser
}

func NewLightningConfigFile() (*LightningConfigFile, error) {
	confPath := path.Join(config.Artifacts, "lightning-config.toml")
	confFile, err := os.Create(confPath)
	if err != nil {
		return nil, err
	}
	return &LightningConfigFile{
		Path: confPath,
		IO:   confFile,
	}, nil
}

func (f *LightningConfigFile) WriteToDisk() string {
	err := f.IO.Close()
	if err != nil {
		log.Warn("failed to close the config file, lightning config might be incorrect", zap.Error(err))
	}
	return f.Path
}
