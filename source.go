package zfs

// Source describes where property is coming from.
type Source string

const (
	// SourceLocal means that property is locally set on given FS.
	SourceLocal Source = "local"

	// SourceDefault means that property is not modified.
	SourceDefault = "default"

	// SourceInherit means that property is inherited from parent FS.
	SourceInherit = "inherit"

	// SourceReceived means that property was received via `zfs receive`.
	SourceReceived = "received"

	// SourceTemporary means that property is set as a result from legacy mount
	// call.
	SourceTemporary = "temporary"

	// SourceNone means that property is internal and has no external source.
	SourceNone = "-"
)
