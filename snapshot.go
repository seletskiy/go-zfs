package zfs

import (
	"errors"
	"io"
	"strings"
)

// See ZFS.NewSnapshot
func NewSnapshot(snapshotPath string) Snapshot {
	return std.NewSnapshot(snapshotPath)
}

// Return Snapshot wrapper without any checks and actualy creation
func (z ZFS) NewSnapshot(snap string) Snapshot {
	buf := strings.Split(snap, "@")
	path := buf[0]
	name := buf[1]
	return Snapshot{zfsEntryBase{z, snap}, NewFS(path), name}
}

type Snapshot struct {
	zfsEntryBase
	FS   FS
	Name string
}

func (s Snapshot) Clone(targetPath string) (FS, error) {
	if s.GetPool() != NewFS(targetPath).GetPool() {
		return FS{}, PoolError
	}

	c := s.runner.Command("zfs", "clone", "-p", s.Path, targetPath)

	_, _, err := c.Output()
	if err != nil {
		return FS{}, err
	}

	return FS{zfsEntryBase{s.runner, targetPath}}, nil
}

func (f FS) Snapshot(name string) (Snapshot, error) {
	snapshotPath := f.Path + "@" + name
	c := f.runner.Command("zfs", "snapshot", snapshotPath)

	_, stderr, err := c.Output()
	if err != nil {
		return Snapshot{}, parseError(err, stderr)
	}

	snap := Snapshot{zfsEntryBase{f.runner, snapshotPath}, f, name}
	return snap, nil
}

func (f FS) ListSnapshots() ([]Snapshot, error) {
	c := f.runner.Command(
		"zfs", "list", "-Hr", "-o", "name", "-t", "snapshot", f.Path,
	)

	stdout, stderr, err := c.Output()
	if err != nil {
		return []Snapshot{}, parseError(err, stderr)
	}

	snapshots := []Snapshot{}
	for _, snap := range strings.Split(strings.TrimSpace(string(stdout)), "\n") {
		if strings.Contains(snap, "@") {
			snapName := strings.Split(snap, "@")[1]
			snapshots = append(snapshots, Snapshot{
				zfsEntryBase{f.runner, snap},
				f,
				snapName,
			})
		} else {
			continue
		}

	}
	return snapshots, nil
}

func notExits(e ZFSEntry) error {
	return errors.New(
		"cannot open '" + e.getPath() + "': dataset does not exist",
	)
}

func (s Snapshot) Send(to ZFSEntry) error {
	rc, stdinPipe, err := to.Receive()
	if err != nil {
		return err
	}

	err = s.SendStream(stdinPipe)
	if err != nil {
		return err
	}

	return parseError(rc.Wait(), nil)
}

func (s Snapshot) SendWithParams(to ZFSEntry) error {
	rc, stdinPipe, err := to.Receive()
	if err != nil {
		return err
	}

	err = s.SendStreamWithParams(stdinPipe)
	if err != nil {
		return err
	}

	return parseError(rc.Wait(), nil)
}

func (s Snapshot) SendIncrementalWithParams(base Snapshot, to ZFSEntry) error {
	rc, stdinPipe, err := to.Receive()
	if err != nil {
		return err
	}

	err = s.SendIncrementalStreamWithParams(base, stdinPipe)
	if err != nil {
		return err
	}

	return parseError(rc.Wait(), nil)
}

func (s Snapshot) SendIncremental(base Snapshot, to ZFSEntry) error {
	rc, stdinPipe, err := to.Receive()
	if err != nil {
		return err
	}

	err = s.SendIncrementalStream(base, stdinPipe)
	if err != nil {
		return err
	}

	return parseError(rc.Wait(), nil)
}

func (s Snapshot) SendStream(dest io.Writer) error {
	if ok, _ := s.Exists(); !ok {
		return notExits(s)
	}

	c := s.runner.Command("zfs", "send", s.Path)

	stdoutPipe, err := c.StdoutPipe()
	if err != nil {
		return err
	}

	if err := c.Start(); err != nil {
		return err
	}

	_, err = io.Copy(dest, stdoutPipe)
	if err != nil {
		return err
	}

	return parseError(c.Wait(), nil)
}

func (s Snapshot) SendStreamWithParams(dest io.Writer) error {
	if ok, _ := s.Exists(); !ok {
		return notExits(s)
	}

	c := s.runner.Command("zfs", "send", "-p", s.Path)

	stdoutPipe, err := c.StdoutPipe()
	if err != nil {
		return err
	}

	if err := c.Start(); err != nil {
		return err
	}

	_, err = io.Copy(dest, stdoutPipe)
	if err != nil {
		return err
	}

	return parseError(c.Wait(), nil)
}

func (s Snapshot) SendIncrementalStream(base Snapshot, dest io.Writer) error {
	if ok, _ := s.Exists(); !ok {
		return notExits(s)
	}
	if ok, _ := base.Exists(); !ok {
		return notExits(base)
	}

	c := s.runner.Command("zfs", "send", "-i", base.Path, s.Path)

	stdoutPipe, err := c.StdoutPipe()
	if err != nil {
		return err
	}

	if err := c.Start(); err != nil {
		return errors.New("error starting send: " + err.Error())
	}

	_, err = io.Copy(dest, stdoutPipe)
	if err != nil {
		return errors.New("error copying to dest: " + err.Error())
	}

	return parseError(c.Wait(), nil)
}

func (s Snapshot) SendIncrementalStreamWithParams(
	base Snapshot,
	dest io.Writer,
) error {
	if ok, _ := s.Exists(); !ok {
		return notExits(s)
	}
	if ok, _ := base.Exists(); !ok {
		return notExits(base)
	}

	c := s.runner.Command("zfs", "send", "-p", "-i", base.Path, s.Path)

	stdoutPipe, err := c.StdoutPipe()
	if err != nil {
		return err
	}

	if err := c.Start(); err != nil {
		return errors.New("error starting send: " + err.Error())
	}

	_, err = io.Copy(dest, stdoutPipe)
	if err != nil {
		return errors.New("error copying to dest: " + err.Error())
	}

	return parseError(c.Wait(), nil)
}

func (s Snapshot) ListClones() ([]FS, error) {
	fss, err := ListFS(s.GetPool())
	if err != nil {
		return []FS{}, err
	}

	clones := []FS{}
	for _, fs := range fss {
		origin, err := fs.GetProperty("origin")
		if err != nil {
			return []FS{}, err
		}

		if origin == s.Path {
			clones = append(clones, fs)
		}
	}

	return clones, nil
}
