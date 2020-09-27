package main

import (
	"github.com/pingcap/log"
	flag "github.com/spf13/pflag"
	components "github.com/yujuncen/brie-bench/workload/components"
	"github.com/yujuncen/brie-bench/workload/utils"
	"github.com/yujuncen/brie-bench/workload/utils/pd"
	"go.uber.org/zap"
)

const Artifacts = "/artifacts"

var (
	component      = flag.String("component", "", "specify the component to test")
	hash           = flag.String("hash", "", "specify the component commit hash")
	repo           = flag.String("repo", "", "specify the repository the bench uses")
	workload       = flag.String("workload", "tpcc1000", "specify the workload")
	debugComponent = flag.Bool("debug-component", false, "component will generate debug level log if enabled")
	disturbance    = flag.Bool("disturbance", false, "enable shuffle-{leader,region,hot-region}-scheduler to simulate extreme environment")
)

func main() {
	flag.Parse()
	cluster := utils.NewCluster()
	if err := utils.DumpCluster(cluster); err != nil {
		log.Warn("failed to dump cluster info", zap.Error(err))
	}
	if err := utils.DumpEnv(); err != nil {
		log.Warn("failed to dump env ", zap.Error(err))
	}
	switch *component {
	case "br":
		br := components.NewBR()
		buildOptions := components.BuildOptions{
			Repository: *repo,
			Hash:       *hash,
		}
		if buildOptions.Repository == "" {
			buildOptions.Repository = "https://github.com/pingcap/br"
		}
		log.Info("build with options", zap.Any("options", buildOptions))
		ibr, err := br.Build(buildOptions)
		utils.Must(err)
		opts := components.BROption{
			Workload:    components.ParseWorkload(*workload),
			LogDir:      Artifacts,
			Cluster:     cluster,
			UseDebugLog: *debugComponent,
		}
		if *disturbance {
			utils.Must(pd.DefaultClient.EnableScheduler([]string{cluster.PdAddr}, pd.Schedulers...))
		}
		log.Info("Run with options", zap.Any("options", opts))
		utils.Must(ibr.Run(opts))
	default:
		log.Panic("Your component isn't supported.\n", zap.String("component", *component))
	}
}
