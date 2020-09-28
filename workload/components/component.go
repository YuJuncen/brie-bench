package component

import (
	"github.com/yujuncen/brie-bench/workload/config"
	"github.com/yujuncen/brie-bench/workload/utils"
)

type BuildOptions struct {
	Repository string
	Hash       string
}

// Component is one of brie components that waiting for testing.
// This interface defines how it could be built.
type Component interface {
	Build(opts BuildOptions) (Binary, error)
	DefaultRepo() string
}

// Binary defines how a built component can be run.
type Binary interface {
	Run(opts interface{}) error
	MakeOptionsWith(conf config.Config, cluster *utils.Cluster) interface{}
}
