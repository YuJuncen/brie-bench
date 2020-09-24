package utils

import (
	"bytes"
	"github.com/pingcap/log"
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

var (
	DropOutput CommandOpt = func(command *exec.Cmd) {
		command.Stdout = nil
		command.Stderr = nil
	}
	SystemOutput CommandOpt = func(command *exec.Cmd) {
		command.Stdout = io.MultiWriter(command.Stdout, os.Stdout)
		command.Stderr = io.MultiWriter(command.Stderr, os.Stderr)
	}
)

func (command *Command) Run() error {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(command.path, command.args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	command.beforeRun(cmd)
	err := cmd.Run()
	log.Debug(cmd.Path, zap.Strings("cmd", cmd.Args),
		zap.String("stdout", stdout.String()), zap.String("stderr", stderr.String()))
	return err
}

func Must(e error) {
	if e != nil {
		log.Panic("meet error %v", zap.Error(e))
	}
}
