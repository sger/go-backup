package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/olekukonko/tablewriter"
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

			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "PATH", "HASH"})

			data := [][]string{}

			var path podule.Path
			b.ForEach(func(k, v []byte) error {

				//fmt.Printf("key=%v, value=%s\n", k, v)
				if err := json.Unmarshal(v, &path); err != nil {
					fmt.Println(err)
				}
				//fmt.Printf("= %v\n", path)

				data = append(data, []string{fmt.Sprintf("%d", path.ID), path.Path, path.Hash})
				table.Append([]string{fmt.Sprintf("%d", path.ID), path.Path, path.Hash})

				return nil
			})

			if len(data) > 0 {
				table.Render()
			} else {
				fmt.Println("no data found")
			}

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

				buf, err := json.Marshal(path)
				if err != nil {
					return err
				}

				err = b.Put(podule.Itob(path.ID), buf)
				if err != nil {
					return err
				}
			}

			return nil
		})
	case "remove":

		if len(args[1:]) == 0 {
			fatalErr = errors.New("specify key id to remove")
			return
		}
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("paths"))

			id, err := strconv.Atoi(args[1:][0])
			if err != nil {
				fmt.Println(err)
			}

			if b.Get(podule.Itob(id)) != nil {
				err = b.Delete(podule.Itob(id))
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("key not exists")
			}
			return err
		})
	}
}
