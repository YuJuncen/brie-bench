package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/yujuncen/brie-bench/workload/config"
	"github.com/yujuncen/brie-bench/workload/utils/metric"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pingcap/errors"
	"github.com/pingcap/log"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"go.uber.org/zap"
)

const (
	ResourcePrefix = "api/cluster/resource/%v"
	ScaleOutPrefix = "api/cluster/scale_out/%v/%v/%v"
	ResultsPrefix  = "api/cluster/workload/%v/result"
)

// ResourceRequestItem ...
type ResourceRequestItem struct {
	ID           uint
	ItemID       uint
	InstanceType string
	RRID         uint
	RID          uint
	Components   string
}

func (r *ResourceRequestItem) hasNum(style string) (num int) {
	components := strings.Split(r.Components, "|")
	for _, component := range components {
		if component == style {
			num++
		}
	}
	return num
}

type WorkloadReport struct {
	CRID      uint
	Data      string
	PlainText *string
}

type BenchContext struct {
	id             string
	name           string
	TidbAddr       string
	PdAddr         string
	PrometheusAddr string
	apiAddr        string
	client         *http.Client

	Report *metric.Report
	Config config.Config
}

func NewCluster(conf config.Config) *BenchContext {
	return &BenchContext{
		id:             os.Getenv("CLUSTER_ID"),
		name:           os.Getenv("CLUSTER_NAME"),
		TidbAddr:       os.Getenv("TIDB_ADDR"),
		PdAddr:         os.Getenv("PD_ADDR"),
		PrometheusAddr: os.Getenv("PROM_ADDR"),
		apiAddr:        os.Getenv("API_SERVER"),
		client:         &http.Client{},

		Config: conf,
		Report: metric.FromConfig(conf),
	}
}

func (c *BenchContext) SetApiServer(apiAddr string) {
	c.apiAddr = apiAddr
}

func (c *BenchContext) SetID(id string) {
	c.id = id
}

func (c *BenchContext) SetName(name string) {
	c.name = name
}

func (c *BenchContext) joinUrl(prefix string) string {
	return c.apiAddr + "/" + prefix
}

func (c *BenchContext) getAllResource() ([]ResourceRequestItem, error) {
	prefix := fmt.Sprintf(ResourcePrefix, c.id)
	url := c.joinUrl(prefix)
	resp, err := doRequest(url, http.MethodGet)
	if err != nil {
		return nil, err
	}
	resources := make([]ResourceRequestItem, 0, 0)
	err = json.Unmarshal([]byte(resp), &resources)
	return resources, err
}

func (c *BenchContext) getAvailableResourceID(component string) (uint, error) {
	resources, err := c.getAllResource()
	if err != nil {
		return 0, errors.New("failed to get all resource")
	}
	// select available
	for _, resource := range resources {
		if resource.hasNum(component) == 0 {
			return resource.ID, nil
		}
	}
	return 0, errors.New("no available resources")
}

func (c *BenchContext) getStoreNum() (num int) {
	resources, err := c.getAllResource()
	if err != nil {
		return 0
	}
	for _, resource := range resources {
		num += resource.hasNum("tikv")
	}
	return num
}

func (c *BenchContext) scaleOut(component string, id uint) error {
	prefix := fmt.Sprintf(ScaleOutPrefix, c.id, id, component)
	url := c.joinUrl(prefix)
	_, err := doRequest(url, http.MethodPost)
	return err
}

func (c *BenchContext) AddStore() error {
	component := "tikv"
	id, err := c.getAvailableResourceID(component)
	if err != nil {
		return err
	}
	return c.scaleOut(component, id)
}

func (c *BenchContext) SendReport(data, plainText string) error {
	prefix := fmt.Sprintf(ResultsPrefix, c.id)
	url := c.joinUrl(prefix)
	return postJSON(url, map[string]interface{}{
		"data":      data,
		"plaintext": plainText,
	})
}

func (c *BenchContext) GetLastReport() (*WorkloadReport, error) {
	prefix := fmt.Sprintf(ResultsPrefix, c.id)
	url := c.joinUrl(prefix)
	resp, err := doRequest(url, http.MethodGet)
	if err != nil {
		return nil, err
	}

	reports := make([]WorkloadReport, 0, 0)
	err = json.Unmarshal([]byte(resp), &reports)
	if err != nil || len(reports) == 0 {
		return nil, err
	}
	return &reports[0], nil
}

func (c *BenchContext) getMetric(query string, t time.Time) (float64, error) {
	client, err := api.NewClient(api.Config{
		Address: c.PrometheusAddr,
	})
	if err != nil {
		log.Error("error creating client", zap.Error(err))
		return 0, err
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := v1api.Query(ctx, query, t)
	if err != nil {
		log.Error("error querying Prometheus", zap.Error(err))
		return 0, err
	}
	if len(warnings) > 0 {
		log.Warn("query has warnings")
	}
	vector := result.(model.Vector)
	if len(vector) >= 1 {
		return float64(vector[0].Value), nil
	}
	return 0, nil
}

func (c *BenchContext) getMatrixMetric(query string, r v1.Range) ([][]float64, error) {
	client, err := api.NewClient(api.Config{
		Address: c.PrometheusAddr,
	})
	if err != nil {
		log.Error("error creating client", zap.Error(err))
		return nil, err
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := v1api.QueryRange(ctx, query, r)
	if err != nil {
		log.Error("error querying Prometheus", zap.Error(err))
		return nil, err
	}
	if len(warnings) > 0 {
		log.Warn("query has warnings")
	}
	matrix := result.(model.Matrix)
	var ret [][]float64
	for _, m := range matrix {
		var r []float64
		for _, v := range m.Values {
			r = append(r, float64(v.Value))
		}
		ret = append(ret, r)
	}
	return ret, nil
}
