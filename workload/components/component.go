package component

import (
	"github.com/yujuncen/brie-bench/workload/config"
	"github.com/yujuncen/brie-bench/workload/utils"
)

type BuildOptions struct {
	Repository string
	Hash       string
}

type Component interface {
	Build(opts BuildOptions) (ComponentBinary, error)
	DefaultRepo() string
}

type ComponentBinary interface {
	Run(opts interface{}) error
	MakeOptionsWith(conf config.Config, cluster *utils.Cluster) interface{}
}
