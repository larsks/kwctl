package main

import (
	"log/slog"
	"os/exec"

	shellquote "github.com/kballard/go-shellquote"
)

// commandRequest represents a request to execute a kwctl command.
type commandRequest struct {
	args        []string
	needsOutput bool
	resultCh    chan commandResult
}

// commandResult represents the result of executing a kwctl command.
type commandResult struct {
	output []byte
	err    error
}

// KwCtl handles execution of kwctl commands with a consistent base command.
// Commands are processed sequentially in FIFO order to prevent conflicts on the serial port.
type KwCtl struct {
	baseArgs  []string
	logger    *slog.Logger
	commandCh chan *commandRequest
}

// NewKwCtl creates a new KwCtl instance by parsing the command string.
// The cmdString should be the base kwctl command (e.g., "kwctl" or "ssh radio kwctl").
func NewKwCtl(cmdString string, logger *slog.Logger) (*KwCtl, error) {
	args, err := shellquote.Split(cmdString)
	if err != nil {
		return nil, err
	}

	k := &KwCtl{
		baseArgs:  args,
		logger:    logger,
		commandCh: make(chan *commandRequest),
	}

	// Start the command processor goroutine
	go k.commandProcessor()

	return k, nil
}

// commandProcessor processes commands sequentially in FIFO order.
func (k *KwCtl) commandProcessor() {
	for req := range k.commandCh {
		allArgs := append(k.baseArgs[1:], req.args...)
		k.logger.Debug("executing kwctl command", "command", k.baseArgs[0], "args", allArgs)

		cmd := exec.Command(k.baseArgs[0], allArgs...)

		if req.needsOutput {
			output, err := cmd.Output()
			req.resultCh <- commandResult{output: output, err: err}
		} else {
			err := cmd.Run()
			req.resultCh <- commandResult{err: err}
		}
	}
}

// Run executes the kwctl command with the given arguments.
// Example: kwctl.Run("vfo", "0") executes "kwctl vfo 0"
// Commands are processed in FIFO order.
func (k *KwCtl) Run(args ...string) error {
	resultCh := make(chan commandResult, 1)
	req := &commandRequest{
		args:        args,
		needsOutput: false,
		resultCh:    resultCh,
	}

	k.commandCh <- req
	result := <-resultCh

	return result.err
}

// RunWithOutput executes the kwctl command and returns its output.
// Example: kwctl.RunWithOutput("status") executes "kwctl status" and returns the output
// Commands are processed in FIFO order.
func (k *KwCtl) RunWithOutput(args ...string) ([]byte, error) {
	resultCh := make(chan commandResult, 1)
	req := &commandRequest{
		args:        args,
		needsOutput: true,
		resultCh:    resultCh,
	}

	k.commandCh <- req
	result := <-resultCh

	return result.output, result.err
}
