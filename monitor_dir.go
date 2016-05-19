package backup

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/sger/go-archiver"
	hashdir "github.com/sger/go-hashdir"
)

// Monitor ...
type Monitor struct {
	Paths       map[string]string
	Destination string
	Archiver    archiver.Archiver
}

// Now ...
func (m *Monitor) Now() (int, error) {
	var counter int
	for path, lastHash := range m.Paths {
		newHash, err := hashdir.Create(path, "md5")
		if err != nil {
			return 0, err
		}
		if newHash != lastHash {
			err := m.Act(path)
			if err != nil {
				return counter, err
			}
			m.Paths[path] = newHash
			counter++

		}
	}
	return counter, nil
}

// Act ...
func (m *Monitor) Act(path string) error {
	fileName := fmt.Sprintf(m.Archiver.Name(), time.Now().UnixNano())
	return m.Archiver.Archive(path, filepath.Join(m.Destination, "", fileName))
}
