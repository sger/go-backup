package main

import "github.com/sger/podule"

func main() {
	m := &podule.Monitor{
		Archiver: podule.ZIP,
		Paths:    make(map[string]string),
	}
}
