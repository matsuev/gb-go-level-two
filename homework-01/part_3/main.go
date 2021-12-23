package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	path = flag.String("path", "./tmp", "Path to create files")
	name = flag.String("name", "file", "Filename")
	num  = flag.Int("num", 1_000_000, "Number of files to create")
)

func main() {
	flag.Parse()
	var err error

	if _, err = os.Stat(*path); os.IsNotExist(err) {
		if err = os.Mkdir(*path, 0777); err != nil {
			log.Fatalln(err)
		}
	}

	for i := 1; i <= *num; i++ {
		if err = createFile(*path, *name, i); err != nil {
			log.Fatalln(err)
		}
	}
}

// createFile function
func createFile(path string, name string, num int) error {
	fh, err := os.Create(fmt.Sprintf("%s/%s-%07d", path, name, num))
	if err != nil {
		return err
	}
	defer fh.Close()

	return nil
}
