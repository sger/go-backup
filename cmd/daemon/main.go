package main

import (
	"encoding/json"
	"fmt"

	"github.com/matryer/filedb"
	"github.com/sger/podule"
)

type path struct {
	Path string
	Hash string
}

func main() {
	m := &podule.Monitor{
		Destination: "test/archive",
		Archiver:    podule.ZIP,
		Paths:       make(map[string]string),
	}

	db, err := filedb.Dial("../../../backupdata")
	if err != nil {
		fmt.Println("Database not found")
		return
	}
	defer db.Close()
	col, err := db.C("paths")
	if err != nil {
		//fmt.Println(err)
	}
	var path path
	col.ForEach(func(_ int, data []byte) bool {
		if err := json.Unmarshal(data, &path); err != nil {
			return true
		}
		m.Paths[path.Path] = path.Hash
		return false
	})
	if len(m.Paths) < 1 {
		fmt.Println("no paths")
	}
}
