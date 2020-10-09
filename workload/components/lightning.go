package component

import (
	"github.com/pingcap/log"
	"github.com/yujuncen/brie-bench/workload/config"
	"github.com/yujuncen/brie-bench/workload/utils"
	"github.com/yujuncen/brie-bench/workload/utils/git"
	"go.uber.org/zap"
	"path"
	"reflect"
	"strings"
)

type Lightning struct{}
type LightningBin struct {
	binary string
}

func (l *LightningBin) MakeOptionsWith(conf config.Config, cluster *utils.Cluster) interface{} {
	var backend LightningBackend
	switch strings.ToLower(conf.Lightning.Backend) {
	case "local":
		backend = Local
	case "tidb":
		backend = TiDB
	default:
		log.Warn("unsupported backend, use local backend to test", zap.String("backend", conf.Lightning.Backend))
		backend = Local
	}
	return LightningOpts{
		Backend: backend,
		Workload: LightningWorkload{
			Name:   conf.Workload,
			Source: conf.WorkloadStorage,
		},
		Cluster: cluster,
		Extra:   conf.ComponentArgs,
	}
}

type LightningWorkload struct {
	Name   string
	Source string
}

type LightningBackend int

const (
	Local LightningBackend = iota
	TiDB
)

type LightningOpts struct {
	Backend  LightningBackend
	Workload LightningWorkload

	Cluster *utils.Cluster

	Extra []string
}

func (l *LightningBin) ImportLocal(opts LightningOpts) error {
	addr, port, err := utils.HostAndPort(opts.Cluster.TidbAddr)
	if err != nil {
		return err
	}
	cliOpts := make([]string, 0)
	cliOpts = append(cliOpts, []string{
		"--backend", "tidb",
		"--tidb-host", addr,
		"--tidb-port", port,
		"--pd-urls", opts.Cluster.PdAddr,
		"-d", opts.Workload.Source,
		"--log-file", path.Join(config.Artifacts, "tidb.log"),
	}...)
	cliOpts = append(cliOpts, opts.Extra...)
	cmd := utils.NewCommand(l.binary, cliOpts...)
	return utils.Bench("import with local backend", cmd.Run)
}

func (l *LightningBin) ImportTiDB(opts LightningOpts) error {
	addr, port, err := utils.HostAndPort(opts.Cluster.TidbAddr)
	if err != nil {
		return err
	}
	cliOpts := make([]string, 0)
	cliOpts = append(cliOpts, []string{
		"--backend", "local",
		"--tidb-host", addr,
		"--tidb-port", port,
		"--pd-urls", opts.Cluster.PdAddr,
		"-d", opts.Workload.Source,
		"--log-file", path.Join(config.Artifacts, "local.log"),
	}...)
	cliOpts = append(cliOpts, opts.Extra...)
	cmd := utils.NewCommand(l.binary, cliOpts...)
	return utils.Bench("import with TiDB backend", cmd.Run)
}

func (l *LightningBin) Run(opts interface{}) error {
	opt, ok := opts.(LightningOpts)
	if !ok {
		log.Error("unexpected config type for lightning", zap.Stringer("type", reflect.TypeOf(opt)))
	}
	switch opt.Backend {
	case Local:
		if err := l.ImportLocal(opt); err != nil {
			return err
		}
	case TiDB:
		if err := l.ImportTiDB(opt); err != nil {
			return err
		}
	}
	return nil
}

func (l Lightning) Build(opts BuildOptions) (Binary, error) {
	repo, err := git.CloneHash(opts.Repository, "/lightning", opts.Hash)
	if err != nil {
		return nil, err
	}
	if err := repo.Make("lightning"); err != nil {
		return nil, err
	}
	return &LightningBin{binary: "/lightning/bin/tidb-lightning"}, nil
}
