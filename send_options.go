package zfs

// SendOptions represents options which can alter Send() process.
type SendOptions struct {
	// Incremental can be used to specify base snapshot name which will be
	// used to create incremental data stream instead of sending full data
	// stream.
	Incremental string

	// IncludeIntermediary will tells Send() that it should send all
	// intermediary snapshots till specified in Incremental option.
	IncludeIntermediary bool

	// ReplicationStream will change Send() so it will send replication stream,
	// which will replicate the specified filesystem, and all descendent file
	// systems, up to the named snapshot. When received, all properties,
	// snapshots, descendent file systems, and clones are preserved.
	ReplicationStream bool

	// LargeBlocks will generate a stream which may contain blocks larger than
	// 128KiB. This flag has no effect if the large_blocks pool feature is
	// disabled, or if the recordsize property of this filesystem has
	// never been set above 128KiB.
	LargeBlocks bool

	// Embedded generate a more compact stream by using WRITE_EMBEDDED records
	// for blocks which are stored more compactly on disk by the embedded_data
	// pool.
	Embedded bool

	// Compressed generate a more compact stream by using compressed WRITE
	// records for blocks which are compressed on disk and in memory (see the
	// compression property for details).
	Compressed bool

	// IncludeProperties will include dataset's properties in the stream.
	IncludeProperties bool
}
