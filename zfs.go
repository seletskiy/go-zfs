package zfs

import "github.com/kovetskiy/runcmd"

type ZFS struct {
	*Runner
}

func NewLocalZFS() (*ZFS, error) {
	runner, err := runcmd.NewLocalRunner()

	return NewZFS(runner), err
}

func NewZFS(runner runcmd.Runner) *ZFS {
	return &ZFS{&Runner{Runner: runner}}
}

func (zfs *ZFS) ListFS() ([]*FS, error) {
	zfs.Command("get", "all", "-H", "-p")
}
