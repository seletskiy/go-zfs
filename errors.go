package zfs

import "fmt"

type (
	// ErrNone means that property has set to none value.
	ErrNone struct {
		Name string
	}

	// ErrEmpty means that property has not set at all.
	ErrEmpty struct {
		Name string
	}

	// ErrNotBool means that property is not `on` or `off` and can't be
	// converted to nil.
	ErrNotBool struct {
		Name  string
		Value string
	}
)

// Error returns string representation of an error.
func (err ErrNone) Error() string {
	return fmt.Sprintf("value of '%s' is 'none'", err.Name)
}

// Error returns string representation of an error.
func (err ErrEmpty) Error() string {
	return fmt.Sprintf("value of '%s' is not set", err.Name)
}

// Error returns string representation of an error.
func (err ErrNotBool) Error() string {
	return fmt.Sprintf(
		"value of '%s' is not 'on' or 'off': '%s'",
		err.Name,
		err.Value,
	)
}
