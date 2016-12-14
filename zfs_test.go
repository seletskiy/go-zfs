package zfs

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/kovetskiy/runcmd"
	"github.com/reconquest/nopio-go"
	"github.com/stretchr/testify/assert"
)

func TestZFS_List_ReturnsFileSystemsList(t *testing.T) {
	test := assert.New(t)

	zfs, err := NewZFS()
	test.NoError(err)

	zfs.SetRunner(&runcmd.MockRunner{
		Stdout: asBytes(
			"zroot/a type filesystem -",
			"zroot/b type filesystem -",
		),
	})

	fs, err := zfs.List("/")
	test.NoError(err)
	test.Len(fs, 2)

	test.True(fs[0].IsFileSystem())
	test.True(fs[1].IsFileSystem())
}

func TestZFS_List_ReturnsFileSystemsWithSize(t *testing.T) {
	test := assert.New(t)

	zfs, err := NewZFS()
	test.NoError(err)

	zfs.SetRunner(&runcmd.MockRunner{
		Stdout: asBytes(
			"zroot/a used 3072 -",
			"zroot/a available 4096 -",
			"zroot/a referenced 2048 -",
		),
	})

	fs, err := zfs.List("/")
	test.NoError(err)
	test.Len(fs, 1)

	size, err := fs[0].GetUsedSize()
	test.NoError(err)
	test.EqualValues(3072, size.AsInt64())
	test.EqualValues("3.0KiB", size.String())

	size, err = fs[0].GetAvailableSize()
	test.NoError(err)
	test.EqualValues(4096, size.AsInt64())
	test.EqualValues("4.0KiB", size.String())

	size, err = fs[0].GetReferencedSize()
	test.NoError(err)
	test.EqualValues(2048, size.AsInt64())
	test.EqualValues("2.0KiB", size.String())
}

func TestZFS_List_ReturnsErrorIfStderrIsNotEmpty(t *testing.T) {
	test := assert.New(t)

	zfs, err := NewZFS()
	test.NoError(err)

	zfs.SetRunner(&runcmd.MockRunner{
		Stderr: asBytes(
			"no such file or directory: /dev/zfs",
		),
	})

	_, err = zfs.List("/")
	test.Error(err)
	test.Contains(err.Error(), "no such file or directory")
}

func TestZFS_Snapshot_ProperlyCallBinary(t *testing.T) {
	test := assert.New(t)

	zfs, err := NewZFS()
	test.NoError(err)

	expectCommand(test, zfs, "zfs", "snapshot", "a@b")

	err = zfs.Snapshot("a", "b")
	test.NoError(err)
}

func TestZFS_Destroy_ProperlyCallBinary(t *testing.T) {
	test := assert.New(t)

	zfs, err := NewZFS()
	test.NoError(err)

	expectCommand(test, zfs, "zfs", "destroy", "blah")

	err = zfs.Destroy("blah")
	test.NoError(err)
}

func TestZFS_Sudo_SetsSudoForNextExecution(t *testing.T) {
	test := assert.New(t)

	zfs, err := NewZFS()
	test.NoError(err)

	zfs.Sudo()

	expectCommand(test, zfs, "sudo", "zfs", "destroy", "blah")

	err = zfs.Destroy("blah")
	test.NoError(err)
}

func TestZFS_DestroyRecursive_ProperlyCallBinary(t *testing.T) {
	test := assert.New(t)

	zfs, err := NewZFS()
	test.NoError(err)

	expectCommand(test, zfs, "zfs", "destroy", "-r", "blah")

	err = zfs.DestroyRecursive("blah", DestroyScopeLocal)
	test.NoError(err)

	expectCommand(test, zfs, "zfs", "destroy", "-R", "blah")

	err = zfs.DestroyRecursive("blah", DestroyScopeGlobal)
	test.NoError(err)
}

func TestZFS_Clone_ProperlyCallsBinary(t *testing.T) {
	test := assert.New(t)

	zfs, err := NewZFS()
	test.NoError(err)

	expectCommand(test, zfs, "zfs", "clone", "a", "b")

	err = zfs.Clone("a", "b")
	test.NoError(err)
}

func TestZFS_CloneWithProperties_PassesProperties(t *testing.T) {
	test := assert.New(t)

	zfs, err := NewZFS()
	test.NoError(err)

	expectCommand(test, zfs, "zfs", "clone", "-o", "x=y,c=d", "a", "b")

	err = zfs.CloneWithProperties(
		"a",
		"b",
		Properties{
			{Name: "x", Value: "y"},
			{Name: "c", Value: "d"},
		},
	)
	test.NoError(err)
}

func TestZFS_Receive_ProperlyCallsBinary(t *testing.T) {
	test := assert.New(t)

	zfs, err := NewZFS()
	test.NoError(err)

	expectCommand(test, zfs, "zfs", "receive", "a")

	err = zfs.Receive("a", nopio.NopReader{}, ReceiveOptions{})
	test.NoError(err)
}

func TestZFS_Send_ReportsProgress(t *testing.T) {
	test := assert.New(t)

	zfs, err := NewZFS()
	test.NoError(err)

	zfs.SetRunner(&runcmd.MockRunner{
		Stderr: asBytes(
			"full\tzroot/a@c\t1075819232",
			"size\t1075819233",
			"15:39:14\t1228816\tzroot/a@c",
			"15:39:15\t2279888\tzroot/a@c",
		),
	})

	sequence := 0

	zfs.SendWithProgress(
		"blah",
		ioutil.Discard,
		SendOptions{},
		func(progress SendProgress) {
			switch sequence {
			case 0:
				test.EqualValues("full", progress.Type)
				test.EqualValues("zroot/a@c", progress.Source)
				test.EqualValues(1075819232, progress.SourceSize)
				test.EqualValues(1075819233, progress.SendSize)
				test.EqualValues(false, progress.Reported)

			case 1:
				test.EqualValues(true, progress.Reported)
				test.EqualValues("15:39:14", progress.Report.Time)
				test.EqualValues(1228816, progress.Report.Size)
				test.EqualValues("zroot/a@c", progress.Report.Snapshot)

			case 2:
				test.EqualValues(true, progress.Reported)
				test.EqualValues("15:39:15", progress.Report.Time)
				test.EqualValues(2279888, progress.Report.Size)
				test.EqualValues("zroot/a@c", progress.Report.Snapshot)
			}

			sequence += 1
		},
	)

	test.EqualValues(3, sequence)
}

func expectCommand(test *assert.Assertions, zfs *ZFS, args ...string) {
	zfs.SetRunner(&runcmd.MockRunner{
		OnCommand: func(worker *runcmd.MockRunnerWorker) {
			test.EqualValues(args, worker.GetArgs())
		},
	})
}

func asBytes(lines ...string) []byte {
	return []byte(strings.Join(lines, "\n"))
}
