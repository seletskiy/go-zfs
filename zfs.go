package zfs

import (
	"bytes"
	"fmt"
	"io"

	"github.com/kovetskiy/runcmd"
	"github.com/reconquest/lexec-go"
	"github.com/reconquest/ser-go"
)

// ZFS is a handle to access various zfs operations, it's not linked to single
// pool and can work with any FS (even on remote FS, if remote runner
// is provided).
type ZFS struct {
	*Runner
}

// NewZFS returns new zfs handle linked to local FS.
func NewZFS() (*ZFS, error) {
	runner, err := runcmd.NewLocalRunner()

	return &ZFS{&Runner{Binary: "zfs", Runner: runner}}, err
}

// Sudo changes state of zfs handle so next execution will be done with
// privileged rights.
func (zfs *ZFS) Sudo() *ZFS {
	zfs.Runner.Sudo = true

	return zfs
}

// SetRunner sets runner which will be used to execute zfs binary.
func (zfs *ZFS) SetRunner(runner runcmd.Runner) *ZFS {
	zfs.Runner.Runner = runner

	return zfs
}

// SetLogger sets logger which will be used to log zfs binary starts, exit
// codes as well as stdout/stderr.
func (zfs *ZFS) SetLogger(logger lexec.Logger) *ZFS {
	zfs.Runner.Logger = logger

	return zfs
}

// List lists all filesystems that are starting with specified prefix.
// Prefix '/' can be used to list all FS.
func (zfs *ZFS) List(prefix string) ([]FS, error) {
	command := zfs.Command("get", "all", "-H", "-p", "-r", prefix)

	stdout, _, err := command.Output()
	if err != nil {
		return nil, err
	}

	reader := bytes.NewReader(stdout)

	var (
		name     string
		property string
		value    string
		source   string
	)

	result := []FS{}

	fs := &FS{}

	for {
		_, err := fmt.Fscanf(
			reader,
			"%s %s %s %s\n",
			&name,
			&property,
			&value,
			&source,
		)

		if name != fs.Name || err == io.EOF {
			if fs.Name != "" {
				result = append(result, *fs)
			}

			fs = &FS{
				Name:       name,
				Properties: map[string]Property{},
			}
		}

		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, ser.Errorf(
				err,
				"error while reading command output: '%s'",
				command.String(),
			)
		}

		fs.Properties[property] = Property{
			Name:   property,
			Value:  value,
			Source: Source(source),
		}
	}

	return result, nil
}

// Snapshot snapshots specified FS.
func (zfs *ZFS) Snapshot(target string, name string) error {
	return zfs.Command("snapshot", target+"@"+name).Execute()
}

// Destroy destroys specified FS.
func (zfs *ZFS) Destroy(target string) error {
	return zfs.Command("destroy", target).Execute()
}

// DestroyRecursive destroys specified FS and all it's descendents. Scope
// of destroy can be controlled with corresponding parameter.
func (zfs *ZFS) DestroyRecursive(target string, scope DestroyScope) error {
	return zfs.Command("destroy", string(scope), target).Execute()
}

// Clone clones specified FS under another name.
func (zfs *ZFS) Clone(source string, target string) error {
	return zfs.CloneWithProperties(source, target, Properties{})
}

// CloneWithProperties is a same as Clone(), but new FS will be given
// specified properties set at creation time.
func (zfs *ZFS) CloneWithProperties(
	source string,
	target string,
	properties Properties,
) error {
	args := []string{"clone"}

	if len(properties) > 0 {
		args = append(args, "-o", properties.String())
	}

	args = append(args, source, target)

	return zfs.Command(args...).Execute()
}

// Receive receives FS into specified name with optional options that controls
// receive process.
func (zfs *ZFS) Receive(
	target string,
	reader io.Reader,
	options ReceiveOptions,
) error {
	args := []string{"receive"}

	if options.ForceRollback {
		args = append(args, "-F")
	}

	if options.Resumable {
		args = append(args, "-s")
	}

	if options.NotMount {
		args = append(args, "-u")
	}

	if options.DiscardFirst {
		args = append(args, "-d")
	}

	if options.DiscardAllButLast {
		args = append(args, "-e")
	}

	if options.Origin != "" {
		args = append(args, "-o", "origin="+options.Origin)
	}

	args = append(args, target)

	command := zfs.Command(args...)
	command.NoLog()
	command.SetStdin(reader)

	return command.Execute()
}

// Send sends specified FS as binary stream which is written in writer with
// options that can control send process.
func (zfs *ZFS) Send(
	source string,
	writer io.Writer,
	options SendOptions,
) error {
	return zfs.SendWithProgress(source, writer, options, nil)
}

// SendWithProgress is a same as Send(), but callback can be specified to
// monitor send progress. Specified callback is guaranteed to be called at
// least once with send size information. Then, callback will be called
// each second and will contains information about current send progress.
func (zfs *ZFS) SendWithProgress(
	source string,
	writer io.Writer,
	options SendOptions,
	callback func(SendProgress),
) error {
	args := []string{"send"}

	if options.Incremental != "" {
		if options.IncludeIntermediary {
			args = append(args, "-I", options.Incremental)
		} else {
			args = append(args, "-i", options.Incremental)
		}
	}

	if options.ReplicationStream {
		args = append(args, "-R")
	}

	if options.LargeBlocks {
		args = append(args, "-L")
	}

	if options.Embedded {
		args = append(args, "-e")
	}

	if options.Compressed {
		args = append(args, "-c")
	}

	if options.IncludeProperties {
		args = append(args, "-p")
	}

	if callback != nil {
		args = append(args, "-v", "-P")
	}

	args = append(args, source)

	command := zfs.Command(args...)
	command.NoLog()
	command.SetStdout(writer)

	stderr, err := command.StderrPipe()
	if err != nil {
		return ser.Errorf(
			err,
			"can't get stderr handle to zfs send %s",
			source,
		)
	}

	err = command.Start()
	if err != nil {
		return ser.Errorf(
			err,
			"can't start zfs send: %s",
			source,
		)
	}

	if callback != nil {
		var progress SendProgress

		_, progress.Error = fmt.Fscanf(
			stderr,
			"%s %s %d\nsize %d\n",
			&progress.Type,
			&progress.Source,
			&progress.SourceSize,
			&progress.SendSize,
		)

		callback(progress)

		for {
			progress.Reported = true

			_, progress.Error = fmt.Fscanf(
				stderr,
				"%s %d %s\n",
				&progress.Report.Time,
				&progress.Report.Size,
				&progress.Report.Snapshot,
			)

			if progress.Error == io.EOF {
				break
			}

			callback(progress)

			if progress.Error != nil {
				break
			}
		}
	}

	return command.Wait()
}
