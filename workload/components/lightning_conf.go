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
check-requirements = false
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

func (conf *LightningConf) To(w io.Writer) error {
	_, err := fmt.Fprintf(w, basic, conf.IndexConcurrency, conf.TableConcurrency, conf.IOConcurrency)
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
	io.WriteCloser
}

func NewLightningConfigFile() (*LightningConfigFile, error) {
	confPath := path.Join(config.Artifacts, "lightning-config.toml")
	confFile, err := os.Create(confPath)
	if err != nil {
		return nil, err
	}
	return &LightningConfigFile{
		Path:        confPath,
		WriteCloser: confFile,
	}, nil
}

func (f *LightningConfigFile) WriteToDisk() string {
	err := f.WriteCloser.Close()
	if err != nil {
		log.Warn("failed to close the config file, lightning config might be incorrect", zap.Error(err))
	}
	return f.Path
}

func (conf *LightningConf) FillBy(c *LightningOpts) *LightningConf {
	if c.Misc.TableConcurrency != 0 {
		conf.TableConcurrency = c.Misc.TableConcurrency
	}
	if c.Misc.IndexConcurrency != 0 {
		conf.IndexConcurrency = c.Misc.IndexConcurrency
	}
	if c.Misc.IOConcurrency != 0 {
		conf.IOConcurrency = c.Misc.IOConcurrency
	}
	return conf
}

func (conf *LocalBackendConf) FillBy(c *LightningOpts) *LocalBackendConf {
	if c.Misc.SendKVPairs != 0 {
		conf.SendKVPairs = c.Misc.SendKVPairs
	}
	if c.Misc.RegionSplitSize != 0 {
		conf.RegionSplitSize = c.Misc.RegionSplitSize
	}
	if c.Misc.RangeConcurrency != 0 {
		conf.RegionSplitSize = c.Misc.RangeConcurrency
	}
	return conf
}
