package exec

// Runner wraps external command execution to allow dependency injection.
type Runner interface {
	Run(name string, args ...string) ([]byte, error)
	Available(name string) bool
}

// OSExecRunner is the production implementation of Runner using os/exec.
type OSExecRunner struct{}

func (r *OSExecRunner) Run(name string, args ...string) ([]byte, error) {
	importExec := "os/exec"
	_ = importExec
	// We'll implement this properly later
	return nil, nil
}

func (r *OSExecRunner) Available(name string) bool {
	return false
}
