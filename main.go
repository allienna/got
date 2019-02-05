package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
)
import "flag"

const gotDirectoryPath = ".got"

var (
	objectsDirectoryPath = fmt.Sprintf("%s/objects", gotDirectoryPath)
	refsDirectoryPath    = fmt.Sprintf("%s/refs", gotDirectoryPath)
	headPath             = fmt.Sprintf("%s/HEAD", gotDirectoryPath)
)

var directoryPaths = []string{
	gotDirectoryPath,
	objectsDirectoryPath,
	fmt.Sprintf("%s/info", objectsDirectoryPath),
	fmt.Sprintf("%s/pack", objectsDirectoryPath),
	refsDirectoryPath,
	fmt.Sprintf("%s/heads", refsDirectoryPath),
	fmt.Sprintf("%s/tags", refsDirectoryPath),
}

func main() {
	flag.NewFlagSet("init", flag.ExitOnError)
	flag.Parse()

	var c int
	var err error

	switch os.Args[1] {
	case "init":
		c, err = initializeRepository()
	default:
		flag.PrintDefaults()
		c = 1
		err = errors.New(fmt.Sprintf("'%s' is not a Got command", os.Args[1]))
	}

	if err != nil {
		fmt.Println(err)
	}

	os.Exit(c)
}

func initializeRepository() (c int, err error) {
	fmt.Println("Init")

	if _, err := os.Stat(gotDirectoryPath); !os.IsNotExist(err) {
		return 1, errors.New("exiting Got project")
	}

	for _, directoryPath := range directoryPaths {
		err = createFolder(directoryPath)
		if err != nil {
			return 1, err
		}
	}

	ref := []byte("ref: refs/heads/master\n")
	err = ioutil.WriteFile(headPath, ref, 0644)
	if err != nil {
		return 1, err
	}

	_, filename, _, _ := runtime.Caller(1)
	gotRepositoryPath := path.Join(path.Dir(filename), gotDirectoryPath)

	fmt.Printf("Initialized empty Got repository in %s\n", gotRepositoryPath)
	return 0, nil
}

func createFolder(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
