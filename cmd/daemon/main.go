package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/boltdb/bolt"
	"github.com/sger/podule"
)

type path struct {
	Path string
	Hash string
}

func main() {
	var fatalErr error
	defer func() {
		if fatalErr != nil {
			flag.PrintDefaults()
			log.Fatalln(fatalErr)
		}
	}()
	var (
		interval = flag.Int("interval", 10, "interval between checks (seconds)")
		archive  = flag.String("archive", "archive", "path to archive location")
		dbpath   = flag.String("db", "./backupdata", "path to database directory")
	)
	flag.Parse()
	fmt.Println(*dbpath)
	fmt.Println(*archive)
	fmt.Println(*interval)

	m := &podule.Monitor{
		Destination: *archive,
		Archiver:    podule.ZIP,
		Paths:       make(map[string]string),
	}
	//../backup/backupdata.db
	db, err := bolt.Open(*dbpath, 0600, nil)
	if err != nil {
		fatalErr = err
		return
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("paths"))
		var path podule.Path
		b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%v, value=%s\n", k, v)
			if err := json.Unmarshal(v, &path); err != nil {
				fatalErr = err
				return fatalErr
			}
			fmt.Println(path.Hash)
			m.Paths[path.Path] = path.Hash
			return nil
		})
		return nil
	})

	if fatalErr != nil {
		return
	}

	if len(m.Paths) < 1 {
		fatalErr = errors.New("no paths - use backup tool to add at least one")
		return
	}

	check(m, db)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-time.After(time.Duration(*interval) * time.Second):
			check(m, db)
		case <-signalChan:
			fmt.Println("stopping...")
			goto stop
		}
	}
stop:
}

func check(m *podule.Monitor, db *bolt.DB) {
	log.Println("Checking...")
	counter, err := m.Now()

	if err != nil {
		log.Fatalln("Failed to backup: ", err)
	}

	if counter > 0 {
		log.Printf("  Archived %d directories\n", counter)
		db.View(func(tx *bolt.Tx) error {

			b := tx.Bucket([]byte("paths"))

			c := b.Cursor()

			var path podule.Path
			for k, v := c.First(); k != nil; k, v = c.Next() {
				fmt.Printf("key=%v, value=%s\n", k, v)
				if err := json.Unmarshal(v, &path); err != nil {
					log.Println("failed to unmarshal data (skipping):", err)
				}
				path.Hash, _ = m.Paths[path.Path]
				fmt.Println("path ", path)
				newData, err := json.Marshal(&path)
				if err != nil {
					log.Println("failed to marshal data (skipping):", err)
				}
				fmt.Println(newData)

				db.Update(func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte("paths"))

					err = b.Put(itob(path.ID), newData)
					if err != nil {
						return err
					}

					return err
				})
			}

			return nil
		})
	} else {
		log.Println("  No changes")
	}
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
