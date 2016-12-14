package zfs

// SendProgress represents progress object which will be reported using
// SendWithProgress() method.
type SendProgress struct {
	// Type is a type of send, which can be full, incremental, etc.
	Type string

	// Source is a name of originally specified FS to send.
	Source string

	// SourceSize is a size of specified FS.
	SourceSize int64

	// SendSize is a estimated size which is how many bytes will be actually
	// sent.
	SendSize int64

	// Reported is true when there is available data about already sent stream.
	// It will be false on the first progress report and true on all further
	// progress reports.
	Reported bool

	// Report represents information about already sent data stream. It will
	// updated every second after Send() was started, meaning, that on the
	// first progress report that structure should not be investigated because
	// there is nothing yet known.
	Report struct {
		// Time is a time token which specifies when report was made (it's
		// local time from source host).
		Time string

		// Size is a size in bytes describing how many bytes of estimated size
		// is sent.
		Size int64

		// Snapshot is a currently transferring snapshot name.
		Snapshot string
	}

	// Error is a any error which occured during obtaining progress report (not
	// send itself).
	Error error
}
