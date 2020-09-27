package main

import (
	"github.com/pingcap/log"
	components "github.com/yujuncen/brie-bench/workload/components"
	"github.com/yujuncen/brie-bench/workload/config"
	"github.com/yujuncen/brie-bench/workload/utils"
	"github.com/yujuncen/brie-bench/workload/utils/pd"
	"go.uber.org/zap"
)

const Artifacts = "/artifacts"

func main() {
	config.Init()
	log.Info("run with", zap.Any("config", config.C))

	cluster := utils.NewCluster()
	if err := utils.DumpCluster(cluster); err != nil {
		log.Warn("failed to dump cluster info", zap.Error(err))
	}
	if err := utils.DumpEnv(); err != nil {
		log.Warn("failed to dump env ", zap.Error(err))
	}
	switch config.C.Component {
	case "br":
		br := components.NewBR()
		buildOptions := components.BuildOptions{
			Repository: config.C.Repo,
			Hash:       config.C.Hash,
		}
		if buildOptions.Repository == "" {
			buildOptions.Repository = "https://github.com/pingcap/br"
		}
		log.Info("build with options", zap.Any("options", buildOptions))
		ibr, err := br.Build(buildOptions)
		utils.Must(err)
		opts := components.BROption{
			Workload:    components.ParseWorkload(config.C.Workload),
			LogDir:      Artifacts,
			Cluster:     cluster,
			UseDebugLog: config.C.DebugComponent,
		}
		if config.C.Disturbance {
			utils.Must(pd.DefaultClient.EnableScheduler([]string{cluster.PdAddr}, pd.Schedulers...))
		}
		log.Info("Run with options", zap.Any("options", opts))
		utils.Must(ibr.Run(opts))
	case "dumpling":
		dumpling := components.NewDumpling()
		buildOptions := components.BuildOptions{
			Repository: config.C.Repo,
			Hash:       config.C.Hash,
		}
		if buildOptions.Repository == "" {
			buildOptions.Repository = "https://github.com/pingcap/dumpling"
		}
		dumpbin, err := dumpling.Build(buildOptions)
		utils.Must(err)
		opts := components.DumplingOpts{
			TargetDir: "/tmp/dumped",
			SplitRows: 0,
			FileType:  config.C.Dumpling.FileType,
			LogPath:   Artifacts,
		}
		utils.Must(dumpbin.Run(opts))

	default:
		log.Panic("Your component isn't supported.\n", zap.String("component", config.C.Component))
	}
}
