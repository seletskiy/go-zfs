package zfs

import "github.com/kovetskiy/runcmd"

type Runner struct {
	runcmd.Runner

	Binary string
	Sudo   bool
}

func (runner *Runner) Command(args ...string) runcmd.CmdWorker {
	name := runner.Binary

	if runner.Sudo {
		args = append([]string{name}, args...)
		name = "sudo"
	}

	return runner.Runner.Command(name, args...)
}
