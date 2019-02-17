package main

import (
	"log"
	"os"
	"path/filepath"
)

func createFileIfNotExists(filename string) string {

	if fileExists(filename) {
		log.Println("use existing file: " + filename)
	} else {
		log.Println("creating new file: " + filename)

		err := os.MkdirAll(filepath.Dir(filename), 0755)
		if err != nil {
			panic(err)
		}

		_, err = os.Create(filename)
		if err != nil {
			panic(err)
		}
	}

	fullpath, err := filepath.Abs(filename)
	if err != nil {
		panic(err)
	}

	return fullpath
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
