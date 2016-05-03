package podule_test

import (
	"strings"
	"testing"

	"github.com/sger/archiver"
	"github.com/sger/podule"
	"github.com/stretchr/testify/require"
)

func TestMonitor(t *testing.T) {
	a := &TestArchiver{}
	m := &podule.Monitor{
		Destination: "test/archive",
		Paths: map[string]string{
			"test/hash1": "abc",
			"test/hash2": "def",
		},
		Archiver: a,
	}

	n, err := m.Now()
	require.NoError(t, err)
	require.Equal(t, 2, n)

	require.Equal(t, 2, len(a.Archives))

	for _, call := range a.Archives {
		require.True(t, strings.HasPrefix(call.Dest, m.Destination))
		require.True(t, strings.HasSuffix(call.Dest, ".zip"))
	}
}

type call struct {
	Src  string
	Dest string
}

type TestArchiver struct {
	Archives []*call
	Restores []*call
}

var _ archiver.Archiver = (*TestArchiver)(nil)

func (a *TestArchiver) Name() string {
	return "%d.zip"
}

func (a *TestArchiver) Archive(src, dest string) error {
	a.Archives = append(a.Archives, &call{Src: src, Dest: dest})
	return nil
}

func (a *TestArchiver) Restore(src, dest string) error {
	a.Restores = append(a.Restores, &call{Src: src, Dest: dest})
	return nil
}
