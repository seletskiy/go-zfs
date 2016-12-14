package zfs

import (
	"strings"
)

// Properties is a list of filesystem properties.
type Properties []Property

// String returns serialized form of list of properties.
func (properties Properties) String() string {
	pairs := []string{}

	for _, property := range properties {
		pairs = append(pairs, property.Name+"="+property.Value)
	}

	return strings.Join(pairs, ",")
}
