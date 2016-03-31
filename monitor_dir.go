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
		fmt.Println("lastHash ", lastHash)
		newHash, err := hashdir.Create(path, "md5")
		fmt.Println("newHash ", newHash)
		if err != nil {
			return 0, err
		}
		if newHash != lastHash {
			fmt.Println("Calling Act")
			err := m.Act(path)
			if err != nil {
				return counter, err
			}
			m.Paths[path] = newHash
			counter++
			fmt.Println("counter ", counter)
		}
	}
	return counter, nil
}

func (m *Monitor) Act(path string) error {
	dirName := filepath.Base(path)
	fmt.Println("dirName ", dirName)
	fileName := fmt.Sprintf(m.Archiver.Name(), time.Now().UnixNano())
	fmt.Println("fileName ", fileName)
	fmt.Println("path ", path)
	fmt.Println("join ", filepath.Join(m.Destination, "", fileName))
	return m.Archiver.Archive(path, filepath.Join(m.Destination, "", fileName))
	//return m.Archiver.Archive(path, filepath.Join(m.Destination, dirName, fileName))
}
