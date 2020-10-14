package metric

import (
	"fmt"
	"github.com/pingcap/log"
	"github.com/yujuncen/brie-bench/workload/config"
	"go.uber.org/zap"
	"io"
	"time"
)

var (
	taskInfoFormat = fmt.Sprintf(`### Task info`+"\n\n```"+
		`
[ %-12s ] %%s
[ %-12s ] %%s
[ %-12s ] %%s
[ %-12s ] %%s
`+"```\n\n", "component", "hash", "repository", "workload")

	comparingTaskInfoFormat = fmt.Sprintf(`### Task info`+"\n\n```"+
		`
[ %-12s ] %%s (vs %%s)
[ %-12s ] %%s (vs %%s)
[ %-12s ] %%s (vs %%s)
[ %-12s ] %%s (vs %%s)
`+"```\n\n", "component", "hash", "repository", "workload")
)

const (
	metricHeader = "### Metrics\n\n| name | value |\n|-|-|"
	metricFormat = "| %-20s | %-10s |\n"

	comparingMetricHeader = "### Metrics\n\n| name | value | diff |\n|-|-|-|"
	comparingMetricFormat = "| %-20s | %-10s | %-10s |\n"
)

type ValueType uint

type Metric struct {
	Name  string `json:"name"`
	Value Value  `json:"value"`
}

type Report struct {
	ComponentName string            `json:"component"`
	ComponentHash string            `json:"hash"`
	ComponentRepo string            `json:"repo"`
	WorkloadName  string            `json:"workload"`
	Metrics       map[string]Metric `json:"metrics"`
}

func FromConfig(conf config.Config) *Report {
	report := &Report{
		ComponentName: conf.Component,
		ComponentHash: conf.Hash,
		ComponentRepo: conf.Repo,
		WorkloadName:  conf.Workload,
		Metrics:       map[string]Metric{},
	}
	if report.ComponentHash == "" {
		report.ComponentHash = "<latest master>"
	}
	if report.ComponentRepo == "" {
		report.ComponentRepo = "<official repository>"
	}
	return report
}

func (report *Report) Add(metric Metric) {
	report.Metrics[metric.Name] = metric
}

func (report *Report) Export(w io.Writer) error {
	_, err := fmt.Fprintf(w, taskInfoFormat, report.ComponentName, report.ComponentHash, report.ComponentRepo, report.WorkloadName)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, metricHeader)
	if err != nil {
		return err
	}
	for _, r := range report.Metrics {
		_, err := fmt.Fprintf(w, metricFormat, r.Name, r.Value.String())
		if err != nil {
			return err
		}
	}
	return err
}

func (report *Report) ExportComparing(w io.Writer, other *Report) error {
	_, err := fmt.Fprintf(w, comparingTaskInfoFormat,
		report.ComponentName, other.ComponentName,
		report.ComponentHash, other.ComponentHash,
		report.ComponentRepo, other.ComponentRepo,
		report.WorkloadName, other.WorkloadName)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(w, comparingMetricHeader)
	if err != nil {
		return err
	}
	for _, r := range report.Metrics {
		othMetric, ok := other.Metrics[r.Name]
		var diffString string
		if !ok {
			diffString = "⚠️ no comparing target"
		} else {
			diffString = formatDiff(r.Value.ComparePercentDiff(othMetric.Value))
		}
		_, err := fmt.Fprintf(w, comparingMetricFormat, r.Name, r.Value.String(), diffString)
		if err != nil {
			return err
		}
	}
	return nil
}

func formatDiff(diff float64) string {
	if diff > 0 {
		return fmt.Sprintf("⬆️ %.2f%%", diff*100)
	}
	if diff < 0 {
		return fmt.Sprintf("⬇️ %.2f%%", diff*100)
	}
	return "🔁 no change"
}

// Bench runs the task, with logging the time cost.
func (report *Report) Bench(name string, task func() error) error {
	start := time.Now()
	defer func() {
		timeCost := time.Since(start)
		log.Info("bench task done", zap.String("name", name), zap.Duration("cost", timeCost))
		report.Add(Duration(timeCost).Named(name))
	}()
	return task()
}
