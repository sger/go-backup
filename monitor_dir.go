package podule

import (
	"fmt"
	"path/filepath"
	"time"

	hashdir "github.com/sger/go-hashdir"
)

type Monitor struct {
	Paths       map[string]string
	Destination string
	Archiver    Archiver
}

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

func (m *Monitor) Act(path string) error {
	//dirName := filepath.Base(path)
	fileName := fmt.Sprintf(m.Archiver.Name(), time.Now().UnixNano())
	return m.Archiver.Archive(path, filepath.Join(m.Destination, "", fileName))
	//return m.Archiver.Archive(path, filepath.Join(m.Destination, dirName, fileName))
}
