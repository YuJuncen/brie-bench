package metric

import (
	"github.com/pingcap/log"
	"go.uber.org/zap"
	"time"
)

type ValueType uint

type Metric struct {
	Name  string `json:"name"`
	Value Value  `json:"value"`
}

type Report struct {
	ComponentName string            `json:"component"`
	WorkloadName  string            `json:"workload"`
	Metrics       map[string]Metric `json:"metrics"`
}

func (report *Report) Add(metric Metric) {
	report.Metrics[metric.Name] = metric
}

// Bench runs the task, with logging the time cost.
func Bench(name string, task func() error) error {
	start := time.Now()
	defer func() {
		log.Info("bench task done", zap.String("name", name), zap.Duration("cost", time.Since(start)))
	}()
	return task()
}
