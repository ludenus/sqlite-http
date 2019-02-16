package main

import(
	"path/filepath"
	"os"
	"log"
)

func createFileIfNotExists(filename string) string {

	if fileExists(filename) {
		log.Println("file exists: " + filename)
	} else {
		log.Println("creating file: " + filename)

		os.MkdirAll(filepath.Dir(filename), 0755)
		os.Create(filename)
	}

	return filename
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
