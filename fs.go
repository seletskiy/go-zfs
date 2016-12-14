package zfs

// Type is a type of FS.
type Type string

const (
	// TypeFileSystem is a common modifiable FS.
	TypeFileSystem Type = "filesystem"

	// TypeSnapshot is a snapshot of common FS.
	TypeSnapshot = "snapshot"

	// TypeClone is a clone of common FS or snapshot.
	TypeClone = "clone"
)

// FS represents single zfs filesystem.
type FS struct {
	// Name is a name of filesystem, typically what is seen in `zfs list` output.
	Name string

	// Properties is a map of all filesystem properties, typically what is seen
	// in the `zfs get` output.
	Properties map[string]Property
}

// GetProperty returns single property by it's name. Even non-existent
// properties can be retrieved, they can be catched by later Property.IsEmpty()
// call.
func (fs *FS) GetProperty(name string) Property {
	if property, ok := fs.Properties[name]; ok {
		return property
	} else {
		return Property{Name: name}
	}
}

// GetUsedSize returns `used` property returned as Size type. It's how much
// disk space is consumed by given FS.
func (fs *FS) GetUsedSize() (Size, error) {
	return fs.GetProperty("used").AsSize()
}

// GetAvailableSize returns `available` property returned as Size type. It's
// how much disk space is available for consuming.
func (fs *FS) GetAvailableSize() (Size, error) {
	return fs.GetProperty("available").AsSize()
}

// GetReferencedSize returns `referenced` property returned as Size type. It's
// how much disk space consumed by FS without underlying snapshots and clones.
func (fs *FS) GetReferencedSize() (Size, error) {
	return fs.GetProperty("referenced").AsSize()
}

// GetMountpoint returns `mountpoint` property.
func (fs *FS) GetMountpoint() (Property, error) {
	return fs.GetProperty("mountpoint"), nil
}

// GetType returns type of given FS.
func (fs *FS) GetType() Type {
	return Type(fs.GetProperty("type").Value)
}

// IsFileSystem returns true if given FS is common filesystem.
func (fs *FS) IsFileSystem() bool {
	return fs.GetType() == TypeFileSystem
}
