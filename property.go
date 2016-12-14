package zfs

import (
	"strconv"
	"time"

	"github.com/reconquest/ser-go"
)

// Property represents single filesystem property.
type Property struct {
	// Name is a property name.
	Name string

	// Value is a property value. Will be empty only if property is not set.
	Value string

	// Source is a source of where property is coming from.
	Source Source
}

// IsReadOnly returns true if property can't be changed (e.g. `used` property).
func (property Property) IsReadOnly() bool {
	return property.Source == SourceNone
}

// IsDefault returns true if property is unchanged and equals default value.
func (property Property) IsDefault() bool {
	return property.Source == SourceDefault
}

// IsNone returns true if option is set to none, which is internal zfs value.
func (property Property) IsNone() bool {
	return property.Value == "none"
}

// IsEmpty returns true if property is not set.
func (property Property) IsEmpty() bool {
	return property.Value == ""
}

// AsBool returns property value as boolean or error if it's not boolean value.
// Typically, boolean values are `on` or `off`.
func (property Property) AsBool() (bool, error) {
	switch property.Value {
	case "on":
		return true, nil
	case "off":
		return false, nil
	default:
		return false, ErrNotBool{property.Name, property.Value}
	}
}

// AsInt64 returns value as int64 value or error if it's can't be converted.
func (property Property) AsInt64() (int64, error) {
	if property.IsNone() {
		return 0, ErrNone{property.Name}
	}

	if property.IsEmpty() {
		return 0, ErrEmpty{property.Name}
	}

	result, err := strconv.ParseInt(property.Value, 10, 64)
	if err != nil {
		return 0, ser.Errorf(
			err,
			"can't convert property value '%s' to 64-bit integer: '%s'",
			property.Name,
			property.Value,
		)
	}

	return result, nil
}

// AsSize returns value as Size type which useful for pretty-printing.
func (property Property) AsSize() (Size, error) {
	size, err := property.AsInt64()
	if err != nil {
		return 0, err
	}

	return Size(size), nil
}

// AsTime returns value as Time.
func (property *Property) AsTime() (time.Time, error) {
	seconds, err := property.AsInt64()
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(seconds, 0), nil
}
