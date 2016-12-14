package zfs

// ReceiveOptions used in conjunction with Receive() method and controls
// behaviour of receive.
type ReceiveOptions struct {
	// ForceRollback will destroy all underlying data and rewrite it with
	// whatever data is received.
	ForceRollback bool

	// Resumable will indicate that if receive is aborted, partially received
	// FS will not be destroyed and receive process can be continued later.
	Resumable bool

	// NotMount will set that receiving FS should not be automatically mounted
	// upon completion.
	NotMount bool

	// DiscardFirst will drop first part of name in specified target FS.
	DiscardFirst bool

	// DiscardAllButLast will drop all but last parts of name in specified
	// target FS.
	DiscardAllButLast bool

	// Origin will override target FS `origin` option.
	Origin string
}
