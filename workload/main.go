package main

import (
	"flag"
	components "github.com/yujuncen/brie-bench/workload/components"

	"github.com/pingcap/log"
	"github.com/yujuncen/brie-bench/workload/utils"
	"go.uber.org/zap"
)

const Artifacts = "/artifacts"

var (
	component = flag.String("component", "", "specify the component to test")
	hash      = flag.String("hash", "", "specify the component commit hash")
	workload  = flag.String("workload", "tpcc1000", "specify the workload")
)

func main() {
	flag.Parse()
	cluster := utils.NewCluster()
	switch *component {
	case "br":
		br := components.NewBR()
		buildOptions := components.BuildOptions{
			Repository: "https://github.com/pingcap/br",
			Hash:       *hash,
		}
		log.Info("build with options", zap.Any("options", buildOptions))
		ibr, err := br.Build(buildOptions)
		utils.Must(err)
		opts := components.BROption{
			Workload: components.ParseWorkload(*workload),
			LogDir:   Artifacts,
			Cluster:  cluster,
		}
		log.Info("Run with options", zap.Any("options", opts))
		utils.Must(ibr.Run(opts))
	default:
		log.Panic("Your component isn't supported.\n", zap.String("component", *component))
	}
}
