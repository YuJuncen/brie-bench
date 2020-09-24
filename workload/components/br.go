package component

import (
	"errors"
	"fmt"
	"github.com/pingcap/log"
	"github.com/yujuncen/brie-bench/workload/utils"
	"go.uber.org/zap"
	"path"
	"time"
)

const S3Args = `access-key=YOURACCESSKEY&secret-access-key=YOURSECRETKEY&endpoint=http://172.16.4.4:30812`
const TempBackupStorage = `local:///tmp/backup`

type BR struct{}
type BRRunner struct {
	binary string
}
type BRWorkload struct {
	backupStorage  string
	restoreStorage string
	name           string
}

type BROption struct {
	Cluster  *utils.Cluster
	LogDir   string
	Workload BRWorkload
}

func (br BRRunner) Run(opts interface{}) error {
	opt, ok := opts.(BROption)
	if !ok {
		return errors.New("bad BR option")
	}
	restoreStart := time.Now()
	restore := utils.NewCommand(br.binary, "restore", "full",
		"--log-file", path.Join(opt.LogDir, "restore.log"),
		"--log-level", "DEBUG",
		"--pd", opt.Cluster.PdAddr,
		"-s", opt.Workload.restoreStorage)
	restore.Opt(utils.DropOutput)
	if err := restore.Run(); err != nil {
		return err
	}
	log.Info("restore done",
		zap.String("workload", opt.Workload.name),
		zap.Stringer("cost", time.Since(restoreStart)))
	backupStart := time.Now()
	backup := utils.NewCommand(br.binary, "backup", "full",
		"--log-file", path.Join(opt.LogDir, "backup.log"),
		"--log-level", "DEBUG",
		"--pd", opt.Cluster.PdAddr,
		"-s", opt.Workload.backupStorage)
	backup.Opt(utils.DropOutput)
	if err := backup.Run(); err != nil {
		return err
	}
	log.Info("backup done",
		zap.String("workload", opt.Workload.name),
		zap.Stringer("cost", time.Since(backupStart)))
	return nil
}

func (B BR) Build(opts BuildOptions) (BuiltComponent, error) {
	if err := utils.NewCommand("git",
		"clone", opts.Repository, "/br").Opt(utils.SystemOutput).Run(); err != nil {
		return nil, err
	}
	if opts.Hash != "" {
		if err := utils.NewCommand("git", "reset", "--hard", opts.Hash).
			Opt(utils.SystemOutput, utils.WorkDir("/br")).
			Run(); err != nil {
			return nil, err
		}
	}
	if err := utils.NewCommand("make", "build").
		Opt(utils.WorkDir("/br"), utils.SystemOutput).
		Run(); err != nil {
		return nil, err
	}
	return BRRunner{"/br/bin/br"}, nil
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
