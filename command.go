package zfs

import (
	"github.com/reconquest/lexec-go"
	"github.com/reconquest/ser-go"
)

// Command represents single command object eligible for execution.
type Command struct {
	*lexec.Execution
}

// Output returns stdout, stderr and error, if program fails to run or returned
// some stderr.
func (command *Command) Output() ([]byte, []byte, error) {
	stdout, stderr, err := command.Execution.Output()
	if err != nil {
		return stdout, stderr, ser.Errorf(
			err,
			"error while execution zfs command '%s'",
			command.String(),
		)
	}

	if len(stderr) != 0 {
		return stdout, stderr, ser.Errorf(
			stderr,
			"command returned non-empty stderr: '%s'",
			command.String(),
		)
	}

	return stdout, stderr, err
}

// Execute is a same as Output(), but ignores any produced stdout.
func (command Command) Execute() error {
	_, _, err := command.Output()

	return err
}
