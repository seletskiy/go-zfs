package zfs

// DestroyScope is a flag which controls how deep destroing will be done.
type DestroyScope string

const (
	// DestroyScopeLocal will set to destroy only direct descendents.
	DestroyScopeLocal DestroyScope = "-r"

	// DestroyScopeGlobal will set to destroy any descendents, even outside
	// of current hierarchy.
	DestroyScopeGlobal = "-R"
)
