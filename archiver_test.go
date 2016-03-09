package podule_test

import (
	"os"
	"testing"

	"github.com/sger/backup"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) {
	os.MkdirAll("test/output", 0777)
}

func teardown(t *testing.T) {
	os.RemoveAll("test/output")
}

func TestZipArchive(t *testing.T) {
	setup(t)
	defer teardown(t)

	err := backup.ZIP.Archive("test/files", "test/output/files.zip")
	require.NoError(t, err)
}

type call struct {
	Src  string
	Dest string
}

type TestArchiver struct {
	Archives []*call
}

var _ backup.Archiver = (*TestArchiver)(nil)

func (a *TestArchiver) DestFmt() string {
	return "%d.zip"
}

func (a *TestArchiver) Archive(src, dest string) error {
	a.Archives = append(a.Archives, &call{Src: src, Dest: dest})
	return nil
}

func (a *TestArchiver) Restore(src, dest string) error {
	return nil
}
