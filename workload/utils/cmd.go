package utils

import (
	"bytes"
	"github.com/pingcap/log"
	"go.uber.org/zap"
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
		command.Stdout = nil
		command.Stderr = nil
	}
	SystemOutput CommandOpt = func(command *exec.Cmd) {
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
	}
)

func (command *Command) Run() error {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(command.path, command.args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	command.beforeRun(cmd)
	log.Info("executing", zap.Stringer("command", cmd))
	err := cmd.Run()
	log.Debug("exec done", zap.String("command", cmd.Path), zap.Strings("args", cmd.Args))
	if cmd.Stderr == &stderr {
		log.Debug("stderr", zap.Stringer("data", &stderr))
	}
	if cmd.Stdout == &stdout {
		log.Debug("stderr", zap.Stringer("data", &stdout))
	}
	return err
}

func Must(e error) {
	if e != nil {
		log.Panic("meet error", zap.Error(e))
	}
}
