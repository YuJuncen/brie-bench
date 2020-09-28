package main

import (
	"errors"
	"github.com/pingcap/log"
	components "github.com/yujuncen/brie-bench/workload/components"
	"github.com/yujuncen/brie-bench/workload/config"
	"github.com/yujuncen/brie-bench/workload/utils"
	"github.com/yujuncen/brie-bench/workload/utils/pd"
	"go.uber.org/zap"
	"time"
)

func startComponent(component components.Component, cluster *utils.Cluster, conf config.Config) error {
	buildOpts := components.BuildOptions{
		Hash:       conf.Hash,
		Repository: conf.Repo,
	}
	if buildOpts.Repository == "" {
		buildOpts.Repository = component.DefaultRepo()
	}
	start := time.Now()
	log.Info("build with options", zap.Any("options", buildOpts))
	bin, err := component.Build(buildOpts)
	if err != nil {
		return err
	}
	runOpts := bin.MakeOptionsWith(conf, cluster)
	log.Info("Run with options", zap.Any("options", runOpts), zap.Duration("build-time-cost", time.Since(start)))
	runStart := time.Now()
	if err := bin.Run(runOpts); err != nil {
		return err
	}
	log.Info("Run ended", zap.Duration("run-time-cost", time.Since(runStart)))
	return nil
}

func parseComponent(name string) (components.Component, error) {
	switch name {
	case "br":
		return components.NewBR(), nil
	case "dumpling":
		return components.NewDumpling(), nil
	default:
		log.Error("Your component isn't supported.", zap.String("component", config.C.Component))
		return nil, errors.New("unsupported component")
	}
}

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
	component, err := parseComponent(config.C.Component)
	utils.Must(err)
	if config.C.Disturbance {
		utils.Must(pd.DefaultClient.EnableScheduler([]string{cluster.PdAddr}, pd.Schedulers...))
	}
	utils.Must(startComponent(component, cluster, config.C))
}
