package utils

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/pingcap/log"
	"github.com/yujuncen/brie-bench/workload/config"
	"go.uber.org/zap"
)

// CommandOpt is an option for a command.
type CommandOpt func(command *exec.Cmd)

// Command is a runnable bash command.
type Command struct {
	path string
	args []string

	beforeRun CommandOpt
}

// NewCommand creates a command.
func NewCommand(path string, args ...string) *Command {
	return &Command{path: path, args: args, beforeRun: func(cmd *exec.Cmd) {}}
}

// Opt apply command options to the command.
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

// WorkDir sets the workdir of the command.
func WorkDir(dir string) CommandOpt {
	return func(command *exec.Cmd) {
		command.Dir = dir
	}
}

// RedirectTo equals `command 2>&1 >$dir`
func RedirectTo(dir string) CommandOpt {
	return func(command *exec.Cmd) {
		file, err := os.Create(dir)
		Must(err)
		command.Stdout = file
		command.Stderr = file
	}
}

var (
	// DropOutput drops the output of the command (Default).
	DropOutput CommandOpt = func(command *exec.Cmd) {
		command.Stdout = nil
		command.Stderr = nil
	}
	// SystemOutput redirect the output to stdin / stdout.
	SystemOutput CommandOpt = func(command *exec.Cmd) {
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
	}
)

// Run runs the command.
func (command *Command) Run() error {
	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	cmd := exec.Command(command.path, command.args...)
	command.beforeRun(cmd)
	if !config.C.DropStdout {
		cmd.Stdout = CombineWriters(cmd.Stdout, stdout)
		cmd.Stderr = CombineWriters(cmd.Stderr, stderr)
	}
	log.Info("executing", zap.Stringer("command", cmd))
	err := cmd.Run()
	if err != nil {
		log.Warn("execute failed", zap.Stringer("command", cmd), zap.Error(err))
		env := new(bytes.Buffer)
		_ = DumpEnvTo(env)
		log.Info("config", zap.Any("config", config.C), zap.Stringer("env", env))
		if !config.C.DropStdout {
			log.Info("stderr", zap.Stringer("data", stderr))
			log.Info("stdout", zap.Stringer("data", stdout))
		}
	}
	return err
}

// Must asserts the error is non-nil.
func Must(e error) {
	if e != nil {
		log.Panic("meet error", zap.Error(e))
	}
}
