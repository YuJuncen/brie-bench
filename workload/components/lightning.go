package component

import (
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/pingcap/log"
	"github.com/yujuncen/brie-bench/workload/config"
	"github.com/yujuncen/brie-bench/workload/utils"
	"github.com/yujuncen/brie-bench/workload/utils/git"
	"go.uber.org/zap"
)

// Lightning is the component tidb-lightning
type Lightning struct{}

// DefaultRepo implements Component
func (l Lightning) DefaultRepo() string {
	return "https://github.com/pingcap/tidb-lightning.git"
}

// LightningBin is a complied lightning
type LightningBin struct {
	binary string
}

// NewLightning creates a new tidb-lightning component.
func NewLightning() Component {
	return Lightning{}
}

// MakeOptionsWith implements Binary.
func (l *LightningBin) MakeOptionsWith(cluster *utils.BenchContext) interface{} {
	conf := cluster.Config
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

// LightningWorkload is a workload spec for lightning
type LightningWorkload struct {
	Name   string
	Source string
}

// LightningBackend is the enum for supported lightning backends.
type LightningBackend int

const (
	// Local is for local backend.
	Local LightningBackend = iota
	// TiDB is for TiDB backend.
	TiDB
)

// LightningOpts is run options for lightning
type LightningOpts struct {
	Backend  LightningBackend
	Workload LightningWorkload
	Misc     config.Lightning

	Cluster *utils.BenchContext

	Extra []string
}

// ImportLocal use lightning local backend to import the workload.
func (l *LightningBin) ImportLocal(opts LightningOpts) error {
	addr, port, err := utils.HostAndPort(opts.Cluster.TidbAddr)
	if err != nil {
		return err
	}

	cliOpts := make([]string, 0)
	conf, err := NewLightningConfigFile()
	if err != nil {
		return nil
	}
	if err := NewLightningConf().FillBy(&opts).To(conf); err != nil {
		return err
	}
	if err := NewLocalBackendConf(path.Join(os.TempDir(), "local-sorting")).
		FillBy(&opts).
		To(conf); err != nil {
		return err
	}

	cliOpts = append(cliOpts, []string{
		"--backend", "local",
		"--tidb-host", addr,
		"--tidb-port", port,
		"--pd-urls", opts.Cluster.PdAddr,
		"-d", opts.Workload.Source,
		"--log-file", path.Join(config.Artifacts, "local.log"),
		"--config", conf.WriteToDisk(),
	}...)

	cliOpts = append(cliOpts, opts.Extra...)
	cmd := utils.NewCommand(l.binary, cliOpts...)
	return opts.Cluster.Report.Bench("import with local backend", cmd.Run)
}

// ImportTiDB use lightning TiDB backend to import the workload.
func (l *LightningBin) ImportTiDB(opts LightningOpts) error {
	addr, port, err := utils.HostAndPort(opts.Cluster.TidbAddr)
	if err != nil {
		return err
	}
	conf, err := NewLightningConfigFile()
	if err != nil {
		return nil
	}
	if err := NewLightningConf().FillBy(&opts).To(conf); err != nil {
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
		"--config", conf.WriteToDisk(),
	}...)
	cliOpts = append(cliOpts, opts.Extra...)
	cmd := utils.NewCommand(l.binary, cliOpts...)
	return opts.Cluster.Report.Bench("import with TiDB backend", cmd.Run)
}

// Run runs lightning.
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

// Build builds the lightning component.
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
