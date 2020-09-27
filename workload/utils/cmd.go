package utils

import (
	"bytes"
	"github.com/pingcap/log"
	"github.com/yujuncen/brie-bench/workload/config"
	"go.uber.org/zap"
	"io"
	"os"
	"os/exec"
)

type CommandOpt func(command *exec.Cmd)

type Command struct {
	path string
	args []string

	beforeRun CommandOpt
}

func NewCommand(path string, args ...string) *Command {
	return &Command{path: path, args: args, beforeRun: func(cmd *exec.Cmd) {}}
}

func (command *Command) Opt(opts ...CommandOpt) *Command {
	oldOpts := command.beforeRun
	command.beforeRun = func(cmd *exec.Cmd) {
		oldOpts(cmd)
		for _, opt := range opts {
			opt(cmd)
		}
	}
	return command
}

func WorkDir(dir string) CommandOpt {
	return func(command *exec.Cmd) {
		command.Dir = dir
	}
}

func RedirectTo(dir string) CommandOpt {
	return func(command *exec.Cmd) {
		file, err := os.Create(dir)
		Must(err)
		command.Stdout = file
		command.Stderr = file
	}
}

var (
	DropOutput CommandOpt = func(command *exec.Cmd) {
		command.Stdout = NopIO
		command.Stderr = NopIO
	}
	SystemOutput CommandOpt = func(command *exec.Cmd) {
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
	}
)

func (command *Command) Run() error {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd := exec.Command(command.path, command.args...)
	cmd.Stdout = NopIO
	cmd.Stderr = NopIO
	command.beforeRun(cmd)
	cmd.Stdout = io.MultiWriter(cmd.Stdout, stdout)
	cmd.Stderr = io.MultiWriter(cmd.Stderr, stderr)
	log.Info("executing", zap.Stringer("command", cmd))
	err := cmd.Run()
	if err != nil {
		log.Warn("execute failed", zap.Stringer("command", cmd), zap.Error(err))
		env := new(bytes.Buffer)
		_ = DumpEnvTo(env)
		log.Info("config", zap.Any("config", config.C), zap.Stringer("env", env))
		log.Info("stderr", zap.Stringer("data", stderr))
		log.Info("stdout", zap.Stringer("data", stdout))
	}
	return err
}

func Must(e error) {
	if e != nil {
		log.Panic("meet error", zap.Error(e))
	}
}
