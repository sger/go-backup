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
)

type path struct {
	ID   int
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
		dbpath = flag.String("db", "./backupdata", "path to database directory")
	)
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fatalErr = errors.New("invalid usage; must specify command")
		return
	}
	fmt.Println(*dbpath)

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
		fmt.Println("Display list")
		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("paths"))

			b.ForEach(func(k, v []byte) error {
				fmt.Printf("key=%v, value=%s\n", k, v)
				return nil
			})
			return nil
		})
	case "add":
		fmt.Println("add")
		if len(args[1:]) == 0 {
			fatalErr = errors.New("must specify path to add")
			return
		}

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("paths"))
			for _, p := range args[1:] {
				path := &path{Path: p, Hash: ""}

				// Generate ID for the user.
				// This returns an error only if the Tx is closed or not writeable.
				// That can't happen in an Update() call so I ignore the error check.
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

				//fmt.Println(path.Path)
				//h := sha256.New()
				//s := fmt.Sprintf("%v", path)
				//sum := h.Sum([]byte(s))
				//fmt.Printf("%s hashes to %x", s, sum)
				fmt.Println(itob(path.ID))
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
