package component

import (
	"errors"
	"fmt"
	"github.com/pingcap/log"
	"github.com/yujuncen/brie-bench/workload/config"
	"github.com/yujuncen/brie-bench/workload/utils"
	"github.com/yujuncen/brie-bench/workload/utils/git"
	"go.uber.org/zap"
	"path"
	"time"
)

const (
	S3Args        = `access-key=YOURACCESSKEY&secret-access-key=YOURSECRETKEY&endpoint=http://172.16.4.4:30812`
	BRDefaultRepo = `https://github.com/pingcap/br`
)

var TempBackupStorage = fmt.Sprintf(`s3://brie/backup?%s`, S3Args)

type BR struct{}

func (B BR) DefaultRepo() string {
	return BRDefaultRepo
}

type BRBin struct {
	binary string
}

func (br BRBin) MakeOptionsWith(conf config.Config, cluster *utils.Cluster) interface{} {
	return BROption{
		Workload:    ParseWorkload(conf.Workload),
		LogDir:      config.Artifacts,
		Cluster:     cluster,
		UseDebugLog: conf.DebugComponent,
		SkipBackup:  conf.BR.SkipBackup,
		Extra:       conf.ComponentArgs,
	}
}

type BRWorkload struct {
	backupStorage  string
	restoreStorage string
	name           string
}

type BRRunType int

type BROption struct {
	Cluster     *utils.Cluster
	LogDir      string
	Workload    BRWorkload
	UseDebugLog bool
	SkipBackup  bool

	Extra []string
}

func (br BRBin) Restore(opt BROption) error {
	restoreStart := time.Now()
	restoreCliOpts := []string{
		"restore", "full",
		"--log-file", path.Join(opt.LogDir, "restore.log"),
		"--pd", opt.Cluster.PdAddr,
		"-s", opt.Workload.restoreStorage,
	}
	if opt.UseDebugLog {
		restoreCliOpts = append(restoreCliOpts, []string{"--log-level", "DEBUG"}...)
	}
	restoreCliOpts = append(restoreCliOpts, opt.Extra...)
	restore := utils.NewCommand(br.binary, restoreCliOpts...)
	restore.Opt(utils.DropOutput)
	if err := restore.Run(); err != nil {
		return err
	}
	log.Info("restore done",
		zap.String("workload", opt.Workload.name),
		zap.Stringer("cost", time.Since(restoreStart)))
	return nil
}

func (br BRBin) Backup(opt BROption) error {
	backupStart := time.Now()
	backupCliOpts := []string{
		"backup", "full",
		"--log-file", path.Join(opt.LogDir, "backup.log"),
		"--pd", opt.Cluster.PdAddr,
		"-s", opt.Workload.backupStorage,
	}
	if opt.UseDebugLog {
		backupCliOpts = append(backupCliOpts, []string{"--log-level", "DEBUG"}...)
	}
	backupCliOpts = append(backupCliOpts, opt.Extra...)
	backup := utils.NewCommand(br.binary, backupCliOpts...)
	backup.Opt(utils.DropOutput)
	if err := backup.Run(); err != nil {
		return err
	}
	log.Info("backup done",
		zap.String("workload", opt.Workload.name),
		zap.Stringer("cost", time.Since(backupStart)))
	return nil
}

func (br BRBin) Run(opts interface{}) error {
	opt, ok := opts.(BROption)
	if !ok {
		return errors.New("bad BR option")
	}
	if err := br.Restore(opt); err != nil {
		return err
	}
	if err := br.Backup(opt); err != nil {
		return err
	}
	return nil
}

func (B BR) Build(opts BuildOptions) (Binary, error) {
	repo, err := git.CloneHash(opts.Repository, "/br", opts.Hash)
	if err != nil {
		return nil, err
	}
	if err := repo.Make("build"); err != nil {
		return nil, err
	}
	return BRBin{"/br/bin/br"}, nil
}

func NewBR() Component {
	return BR{}
}

func TPCCWorkload() BRWorkload {
	return BRWorkload{
		backupStorage:  TempBackupStorage,
		restoreStorage: fmt.Sprintf("s3://mybucket/tpcc1000?%s", S3Args),
		name:           "TPCC-1000",
	}
}

func YCSBWorkload() BRWorkload {
	return BRWorkload{
		backupStorage:  TempBackupStorage,
		restoreStorage: fmt.Sprintf("s3://mybucket/ycsb?%s", S3Args),
		name:           "YCSB",
	}
}

func ParseWorkload(workload string) BRWorkload {
	switch workload {
	case "tpcc1000":
		return TPCCWorkload()
	case "ycsb":
		return YCSBWorkload()
	default:
		return TPCCWorkload()
	}
}
