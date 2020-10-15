package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"time"

	"github.com/pingcap/log"
	components "github.com/yujuncen/brie-bench/workload/components"
	"github.com/yujuncen/brie-bench/workload/config"
	"github.com/yujuncen/brie-bench/workload/utils"
	"github.com/yujuncen/brie-bench/workload/utils/metric"
	"github.com/yujuncen/brie-bench/workload/utils/pd"
	"go.uber.org/zap"
)

func startComponent(component components.Component, cluster *utils.BenchContext, conf config.Config) error {
	buildOpts := components.BuildOptions{
		Hash:       conf.Hash,
		Repository: conf.Repo,
	}
	report, err := cluster.GetLastReport()
	if err != nil {
		log.Warn("failed to get last report, assuming std test case", utils.ShortError(err))
		err = nil
	}
	if buildOpts.Repository == "" || (report != nil && report.Data != "") /* when in PR benching mode */ {
		buildOpts.Repository = component.DefaultRepo()
	}
	start := time.Now()
	log.Info("build with options", zap.Any("options", buildOpts))
	bin, err := component.Build(buildOpts)
	if err != nil {
		return err
	}
	runOpts := bin.MakeOptionsWith(cluster)
	log.Info("Run with options", zap.Any("options", runOpts), zap.Duration("build-time-cost", time.Since(start)))
	runStart := time.Now()
	if err := bin.Run(runOpts); err != nil {
		return err
	}
	log.Info("Run ended", zap.Duration("run-time-cost", time.Since(runStart)))

	return saveAndUploadReport(cluster, report)
}

func saveAndUploadReport(cluster *utils.BenchContext, report *utils.WorkloadReport) error {
	reportFile, err := os.Create(config.Report)
	if err != nil {
		return err
	}
	defer func() { _ = reportFile.Close() }()
	if report != nil && report.Data != "" {
		var lastReport metric.Report
		if err := json.Unmarshal([]byte(report.Data), &lastReport); err != nil {
			return err
		}
		if err := cluster.Report.ExportComparing(reportFile, &lastReport); err != nil {
			return err
		}
	} else {
		jsonReport, err := json.Marshal(cluster.Report)
		if err != nil {
			return err
		}
		textualReport := bytes.NewBuffer(nil)
		if err := cluster.Report.Export(io.MultiWriter(textualReport, reportFile)); err != nil {
			return err
		}
		if err := cluster.SendReport(string(jsonReport), textualReport.String()); err != nil {
			return err
		}
	}
	return nil
}

func parseComponent(name string) (components.Component, error) {
	switch name {
	case "br":
		return components.NewBR(), nil
	case "dumpling":
		return components.NewDumpling(), nil
	case "lightning":
		return components.NewLightning(), nil
	default:
		log.Error("Your component isn't supported.", zap.String("component", config.C.Component))
		return nil, errors.New("unsupported component")
	}
}

func main() {
	config.Init()
	log.Info("run with", zap.Any("config", config.C))

	cluster := utils.NewCluster(config.C)
	if err := utils.DumpCluster(cluster); err != nil {
		log.Warn("failed to dump cluster info", zap.Error(err))
	}
	if err := utils.DumpEnv(); err != nil {
		log.Warn("failed to dump env", zap.Error(err))
	}
	component, err := parseComponent(config.C.Component)
	utils.Must(err)
	if config.C.Disturbance {
		utils.Must(pd.DefaultClient.EnableScheduler([]string{cluster.PdAddr}, pd.Schedulers...))
	}
	utils.Must(startComponent(component, cluster, config.C))
}
