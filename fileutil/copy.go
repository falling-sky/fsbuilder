package fileutil

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func copyHelper(source string, dest string, fn func(string) ([]string, error)) {
	files, err := fn(source)
	if err != nil {
		log.Fatal(err)
	}
	seen := make(map[string]bool)
	for _, f := range files {
		if strings.HasSuffix(f, "~") {
			continue // Skip editor backups
		}
//		log.Printf("copy %s/%s to %s/%s\n", source, f, dest, f)

		// Read the file.
		b, e := ioutil.ReadFile(source + "/" + f)
		if e != nil {
			log.Fatal(err)
		}

		// Create directory, if needed.
		dir := filepath.Dir(dest + "/" + f)
		if _, ok := seen[dir]; ok == false {
			seen[dir] = true
			os.MkdirAll(dir, 0755)
		}

		// Write the file.
		e = ioutil.WriteFile(dest+"/"+f, b, 0644)
		if e != nil {
			log.Fatal(err)
		}
	}
}

func CopyFiles(source string, dest string) {
	log.Printf("copyFiles(%s,%s)\n", source, dest)
	copyHelper(source, dest, FilesInDirNotRecursive)
}

func CopyFilesAll(source string, dest string) {
	log.Printf("copyFiles(%s,%s)\n", source, dest)
	copyHelper(source, dest, FilesInDirRecursive)
}
