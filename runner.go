package zfs

import "github.com/kovetskiy/runcmd"
import "github.com/reconquest/lexec-go"

// Runner represents object which will be used to execute zfs binaries.
type Runner struct {
	// runcmd.Runner is a underlying remote, local or test mock runner.
	runcmd.Runner

	// Logger is a logger function which will be called to zfs calls stdout
	// and stderr.
	Logger lexec.Logger

	// Binary specifies zfs binary name (`zfs` by default).
	Binary string

	// Sudo is a flag which controls that next execution will be run under
	// privileged rights.
	Sudo bool
}

// Command returns object which is suitable for later execution. Binary name
// is prepended automatically, so args should not contain `zfs`.
func (runner *Runner) Command(args ...string) Command {
	name := runner.Binary

	if runner.Sudo {
		args = append([]string{name}, args...)
		name = "sudo"
		runner.Sudo = false
	}

	return Command{
		lexec.New(runner.Logger, runner.Runner.Command(name, args...)),
	}
}
