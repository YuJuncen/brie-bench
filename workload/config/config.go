package config

import (
	flag "github.com/spf13/pflag"
)

type Config struct {
	Component      string
	Hash           string
	Repo           string
	Workload       string
	DebugComponent bool
	Disturbance    bool

	Dumpling Dumpling
}

type Dumpling struct {
	FileType string
}

// C is the global config object
var C Config

func Init() {
	flag.StringVar(&C.Component, "component", "", "specify the component to test")
	flag.StringVar(&C.Hash, "hash", "", "specify the component commit hash")
	flag.StringVar(&C.Repo, "repo", "", "specify the repository the bench uses")
	flag.StringVar(&C.Workload, "workload", "tpcc1000", "specify the workload")
	flag.BoolVar(&C.DebugComponent, "debug-component", false, "component will generate debug level log if enabled")
	flag.BoolVar(&C.Disturbance, "disturbance", false, "enable shuffle-{leader,region,hot-region}-scheduler to simulate extreme environment")

	flag.StringVar(&C.Dumpling.FileType, "", "sql", "the file type of dumpling")
	flag.Parse()
}