package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/sger/podule"
)

func main() {
	var fatalErr error
	defer func() {
		if fatalErr != nil {
			flag.PrintDefaults()
			log.Fatalln(fatalErr)
		}
	}()
	var (
		dbpath = flag.String("db", "./backupdata", "path to database directory")
	)
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fatalErr = errors.New("invalid usage; must specify command")
		return
	}

	fmt.Println(args)
	fmt.Println(dbpath)

	db, err := bolt.Open(*dbpath, 0600, nil)
	if err != nil {
		fatalErr = err
		return
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("paths"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		fatalErr = err
		return
	}
	switch strings.ToLower(args[0]) {
	case "list":
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("paths"))

			var path podule.Path
			b.ForEach(func(k, v []byte) error {
				//fmt.Printf("key=%v, value=%s\n", k, v)
				if err := json.Unmarshal(v, &path); err != nil {
					fmt.Println(err)
				}
				fmt.Printf("= %v\n", path)
				return nil
			})
			return nil
		})
	case "add":
		if len(args[1:]) == 0 {
			fatalErr = errors.New("must specify path to add")
			return
		}

		db.Update(func(tx *bolt.Tx) error {

			b := tx.Bucket([]byte("paths"))

			for _, p := range args[1:] {

				path := &podule.Path{Path: p, Hash: "Not yet archived"}

				id, err := b.NextSequence()
				if err != nil {
					return err
				}

				path.ID = int(id)

				// Marshal user data into bytes.
				buf, err := json.Marshal(path)
				if err != nil {
					return err
				}

				err = b.Put(itob(path.ID), buf)
				if err != nil {
					return err
				}
			}

			return nil
		})
	case "remove":
		fmt.Println("remove")
	}
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
