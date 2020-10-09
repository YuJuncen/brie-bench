package component

import (
	"errors"
	"fmt"
	"github.com/pingcap/log"
	"github.com/yujuncen/brie-bench/workload/config"
	"github.com/yujuncen/brie-bench/workload/utils"
	"github.com/yujuncen/brie-bench/workload/utils/git"
	"github.com/yujuncen/brie-bench/workload/utils/storage"
	"go.uber.org/zap"
	"path"
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
	opt := BROption{
		Workload: BRWorkload{
			Name:              conf.Workload,
			BackupStorageURL:  TempBackupStorage,
			RestoreStorageURL: conf.WorkloadStorage,
		},
		LogDir:      config.Artifacts,
		Cluster:     cluster,
		UseDebugLog: conf.DebugComponent,
		SkipBackup:  conf.BR.SkipBackup,
		Extra:       conf.ComponentArgs,
	}
	if conf.TemporaryStorage != "" {
		log.Info("use other temporary storage", zap.String("url", conf.TemporaryStorage))
		opt.Workload.BackupStorageURL = conf.TemporaryStorage
	}
	return opt
}

type BRWorkload struct {
	BackupStorageURL  string
	RestoreStorageURL string
	Name              string
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
	restoreCliOpts := []string{
		"restore", "full",
		"--log-file", path.Join(opt.LogDir, "restore.log"),
		"--pd", opt.Cluster.PdAddr,
		"-s", opt.Workload.RestoreStorageURL,
	}
	if opt.UseDebugLog {
		restoreCliOpts = append(restoreCliOpts, []string{"--log-level", "DEBUG"}...)
	}
	restoreCliOpts = append(restoreCliOpts, opt.Extra...)
	restore := utils.NewCommand(br.binary, restoreCliOpts...)
	restore.Opt(utils.DropOutput)
	return utils.Bench(fmt.Sprintf("restore %s", opt.Workload.Name), restore.Run)
}

func (br BRBin) Backup(opt BROption, storage *storage.TempS3Storage) error {
	defer func() {
		if err := storage.Cleanup(); err != nil {
			log.Info("failed to cleanup backup", zap.Error(err), zap.String("storage", storage.Raw))
		}
	}()
	backupCliOpts := []string{
		"backup", "full",
		"--log-file", path.Join(opt.LogDir, "backup.log"),
		"--pd", opt.Cluster.PdAddr,
		"-s", storage.Raw,
	}
	if opt.UseDebugLog {
		backupCliOpts = append(backupCliOpts, []string{"--log-level", "DEBUG"}...)
	}
	backupCliOpts = append(backupCliOpts, opt.Extra...)
	backup := utils.NewCommand(br.binary, backupCliOpts...)
	backup.Opt(utils.DropOutput)
	return utils.Bench(fmt.Sprintf("backup %s", opt.Workload.Name), backup.Run)
}

func (br BRBin) Run(opts interface{}) error {
	opt, ok := opts.(BROption)
	if !ok {
		return errors.New("bad BR option")
	}
	if err := br.Restore(opt); err != nil {
		return err
	}
	s, err := storage.ConnectToS3(opt.Workload.BackupStorageURL)
	if err != nil {
		return err
	}
	if err := br.Backup(opt, s); err != nil {
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
