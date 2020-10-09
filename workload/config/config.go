package config

import (
	flag "github.com/spf13/pflag"
)

type Config struct {
	Component        string
	Hash             string
	Repo             string
	Workload         string
	WorkloadStorage  string
	DebugComponent   bool
	Disturbance      bool
	TemporaryStorage string

	ComponentArgs []string

	Dumpling  Dumpling
	BR        BR
	Lightning Lightning
}

type BR struct {
	SkipBackup bool
}

type Dumpling struct {
	SkipCSV bool
	SkipSQL bool
}

type Lightning struct {
	Backend string
	Source  string
}

// C is the global config object
var C Config

func Init() {
	flag.StringVar(&C.Component, "component", "", "specify the component to test")
	flag.StringVar(&C.Hash, "hash", "", "specify the component commit hash")
	flag.StringVar(&C.Repo, "repo", "", "specify the repository the bench uses")
	flag.StringVar(&C.Workload, "workload-name", "", "specify the workload name")
	flag.StringVar(&C.WorkloadStorage, "workload-storage", "", "(with br syntax) specify the storage for workload")
	flag.BoolVar(&C.DebugComponent, "debug-component", false, "component will generate debug level log if enabled")
	flag.BoolVar(&C.Disturbance, "disturbance", false, "enable shuffle-{leader,region,hot-region}-scheduler to simulate extreme environment")
	flag.StringSliceVar(&C.ComponentArgs, "cargs", []string{}, "(unsafe) pass extra argument to the component, may conflict with args provided by the framework")
	flag.StringVar(&C.TemporaryStorage, "temp-storage", "", "(with br syntax) specify the storage where the intermedia data stores")

	flag.BoolVar(&C.Dumpling.SkipCSV, "dumpling.skip-csv", false, "skip dumpling to csv step in dumpling benching")
	flag.BoolVar(&C.Dumpling.SkipSQL, "dumpling.skip-sql", false, "skip dumpling to sql step in dumpling benching")
	flag.BoolVar(&C.BR.SkipBackup, "br.skip-backup", false, "skip the backup step of br benching")
	flag.StringVar(&C.Lightning.Backend, "lightning.backend", "local", "the backend that lightning uses for benching")
	flag.Parse()
}
