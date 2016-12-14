package zfs

import (
	"strconv"
)

// SizeScale is a object which describes how to represent Size in
// human-readable form with prefix.
type SizeScale struct {
	// Divizor is a base which will be used to divide incoming size in bytes.
	Divizor int

	// Suffixes is a slice of suffixes in increasing order.
	Suffixes []string
}

var (
	// SizeScaleBinary is a typical binary scaler (e.g. powers of 1024).
	SizeScaleBinary = SizeScale{
		Divizor: 1024,
		Suffixes: []string{
			"B",
			"KiB",
			"MiB",
			"GiB",
			"TiB",
			"PiB",
			"EiB",
			"ZiB",
		},
	}
)

// Size represents byte size with human-readable form.
type Size int64

// Format returns formatted size given specified precision and scaling.
func (size Size) Format(precision int, scale SizeScale) string {
	var (
		power = 0
		value = float64(size)
	)

	for value >= float64(scale.Divizor) && power+1 < len(scale.Suffixes) {
		power += 1
		value /= float64(scale.Divizor)
	}

	return strconv.FormatFloat(value, 'f', precision, 64) +
		scale.Suffixes[power]
}

// String returns commonly formatted size with one sign after period and using
// binary-powered prefixes like MiB.
func (size *Size) String() string {
	return size.Format(1, SizeScaleBinary)
}

// AsInt64 converts size back to 64-bit integer.
func (size Size) AsInt64() int64 {
	return int64(size)
}
