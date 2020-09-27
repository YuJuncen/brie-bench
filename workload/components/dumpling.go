package component

import (
	"errors"
	"github.com/pingcap/log"
	"github.com/yujuncen/brie-bench/workload/utils"
	"github.com/yujuncen/brie-bench/workload/utils/git"
	"go.uber.org/zap"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

type Dumpling struct{}

func NewDumpling() Component {
	return Dumpling{}
}

func (d Dumpling) Build(opts BuildOptions) (BuiltComponent, error) {
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

func (d *DumplingBin) Run(opts interface{}) error {
	opt, ok := opts.(DumplingOpts)
	if !ok {
		return errors.New("dumpling running with incompatible opt")
	}
	begin := time.Now()
	binOpts := []string{
		"--output", opt.TargetDir,
		"--filetype", opt.FileType,
	}
	if opt.SplitRows > 0 {
		binOpts = append(binOpts, []string{"--rows", strconv.Itoa(opt.SplitRows)}...)
	}
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
	log.Info("dumpling done", zap.Duration("cost", time.Since(begin)))
	return nil
}

type DumplingFileType int

type DumplingOpts struct {
	TargetDir string
	SplitRows int
	FileType  string
	LogPath   string
}
