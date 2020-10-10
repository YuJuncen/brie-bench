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
	// Build builds the component, return a runnable binary.
	Build(opts BuildOptions) (Binary, error)
	// DefaultRepo returns the default git repository for this component.
	DefaultRepo() string
}

// Binary defines how a built component can be run.
type Binary interface {
	// Run runs the component with specified option. The type of option should as same as MakeOptionsWith returns.
	Run(opts interface{}) error
	// MakeOptionsWith extract config value this component needs.
	MakeOptionsWith(conf config.Config, cluster *utils.Cluster) interface{}
}
