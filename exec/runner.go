package exec

import (
	importExec "os/exec"
)

type Runner interface {
	Run(name string, args ...string) ([]byte, error)
	Available(name string) bool
}

type OSExecRunner struct{}

func (r *OSExecRunner) Run(name string, args ...string) ([]byte, error) {
	cmd := importExec.Command(name, args...)
	return cmd.Output()
}

func (r *OSExecRunner) Available(name string) bool {
	_, err := importExec.LookPath(name)
	return err == nil
}
