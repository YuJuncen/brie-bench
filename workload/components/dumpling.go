package component

import (
	"errors"
	"github.com/pingcap/log"
	"github.com/yujuncen/brie-bench/workload/config"
	"github.com/yujuncen/brie-bench/workload/utils"
	"github.com/yujuncen/brie-bench/workload/utils/git"
	"go.uber.org/zap"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	DumplingDefaultRepo = `https://github.com/pingcap/dumpling`
)

type Dumpling struct{}

func (d Dumpling) DefaultRepo() string {
	return DumplingDefaultRepo
}

func NewDumpling() Component {
	return Dumpling{}
}

func (d Dumpling) Build(opts BuildOptions) (Binary, error) {
	repo, err := git.CloneHash(opts.Repository, "/dumpling", opts.Hash)
	if err != nil {
		return nil, err
	}
	if err := repo.Make("build"); err != nil {
		return nil, err
	}
	return &DumplingBin{
		binary: "/dumpling/bin/dumpling",
	}, nil
}

type DumplingBin struct {
	binary string
}

func (d *DumplingBin) MakeOptionsWith(conf config.Config, cluster *utils.Cluster) interface{} {
	return DumplingOpts{
		TargetDir: "/tmp/dumped",
		LogPath:   config.Artifacts,
		Workload:  conf.Workload,
		Cluster:   cluster,

		SkipSQL: conf.Dumpling.SkipSQL,
		SkipCSV: conf.Dumpling.SkipCSV,

		Extra: conf.ComponentArgs,
	}
}

func (d *DumplingBin) Dump(opt DumplingOpts, fileType string) error {
	begin := time.Now()
	host, port, err := utils.HostAndPort(opt.Cluster.TidbAddr)
	if err != nil {
		return err
	}
	binOpts := []string{
		"--output", opt.TargetDir,
		"--filetype", fileType,
		"--host", host,
		"--port", port,
	}
	binOpts = append(binOpts, opt.Extra...)
	if err := utils.NewCommand(d.binary, binOpts...).
		Opt(utils.RedirectTo(path.Join(opt.LogPath, "dumpling.log"))).
		Run(); err != nil {
		return err
	}
	if err := filepath.Walk(opt.TargetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		log.Info("file dumped", zap.String("name", info.Name()), zap.Stringer("size", utils.Size(info.Size())))
		return nil
	}); err != nil {
		return err
	}
	log.Info("dumpling done", zap.Duration("cost", time.Since(begin)),
		zap.String("workload", opt.Workload),
		zap.String("filetype", fileType))
	return nil
}

func (d *DumplingBin) Cleanup(opts DumplingOpts) error {
	log.Info("removing dumped files", zap.String("path", opts.TargetDir))
	return os.RemoveAll(opts.TargetDir)
}

func (d *DumplingBin) Run(opts interface{}) error {
	opt, ok := opts.(DumplingOpts)
	if !ok {
		return errors.New("dumpling running with incompatible opt")
	}
	if !opt.SkipCSV {
		if err := d.Dump(opt, "csv"); err != nil {
			return err
		}
	}
	if err := d.Cleanup(opt); err != nil {
		log.Warn("failed to cleanup dumpling result", zap.Error(err))
		err = nil
	}
	if !opt.SkipSQL {
		if err := d.Dump(opt, "sql"); err != nil {
			return err
		}
	}
	if err := d.Cleanup(opt); err != nil {
		log.Warn("failed to cleanup dumpling result", zap.Error(err))
		err = nil
	}
	return nil
}

type DumplingOpts struct {
	TargetDir string
	LogPath   string
	Workload  string
	SkipSQL   bool
	SkipCSV   bool

	Cluster *utils.Cluster

	Extra []string
}
