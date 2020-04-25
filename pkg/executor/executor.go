package executor

import (
	"context"
	"errors"
	"os/exec"
	"time"

	"github.com/gojavacoder/go-server-monitoring-agent/pkg/util"
)

const EMPTY_COMMAND = "Can't execute empty command"
const COMMAND_TIMEOUT = "Command timed out"
const ERROR_WHILE_EXECUTING_COMMAND = "Error while executing command"

type Executor struct {
	command       string
	commandArgs   []string
	timeout       int
	output        string
	err           error
	isSuccess     bool
	executionTime int64
}

func (e *Executor) SetCommand(com string) {
	e.command = com
}

func (e *Executor) GetCommand() string {
	return e.command
}

func (e *Executor) SetCommandArgs(comArgs []string) {
	e.commandArgs = comArgs
}

func (e *Executor) GetCommandArgs() []string {
	return e.commandArgs
}

func (e *Executor) SetTimeout(timeout int) {
	e.timeout = timeout
}

func (e *Executor) GetTimeout() int {
	return e.timeout
}

func (e *Executor) GetExecutionTime() int64 {
	return e.executionTime
}

func (e *Executor) GetOutput() string {
	return e.output
}

func (e *Executor) IsSuccess() bool {
	return e.isSuccess
}

func (e *Executor) GetError() error {
	return e.err
}

func (e *Executor) validate() bool {
	if e.command == "" {
		e.err = errors.New(EMPTY_COMMAND)
		return false
	} else {
		return true
	}
}

func (e *Executor) setIsSuccess(isSuccess bool) {
	e.isSuccess = isSuccess
	e.executionTime = util.NowAsUnixMilli()
}

func (e *Executor) Execute() bool {
	var err error
	var out []byte

	if !e.validate() {
		e.setIsSuccess(false)
		return e.isSuccess
	}

	ctx, cancel := context.WithTimeout(context.Background(), (time.Duration(e.timeout) * time.Second))
	defer cancel() // The cancel should be deferred so resources are cleaned up

	cmd := exec.CommandContext(ctx, e.command, e.commandArgs...)

	out, err = cmd.Output()

	if ctx.Err() == context.DeadlineExceeded {
		e.output = COMMAND_TIMEOUT
	}

	if err != nil {
		e.err = err
		e.output = ERROR_WHILE_EXECUTING_COMMAND
		e.setIsSuccess(false)
	} else {
		e.output = string(out[:])
		e.setIsSuccess(true)
	}

	return e.isSuccess
}
