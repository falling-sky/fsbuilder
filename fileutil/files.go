package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FilesInDirRecursive returns a list of all files below
// this root directory.  The root directory name is removed
// from the files; the returned names are relative to the root.
func FilesInDirRecursive(root string) ([]string, error) {
	found := []string{}

	fi, err := os.Stat(root)
	switch {
	case err != nil:
		return found, err
	case fi.IsDir():
		if strings.HasSuffix(root, "/") == false {
			root = root + "/"
		}
	default:
		return found, fmt.Errorf("%v: Not a directory", root)
	}

	// Create a callback function
	walker := func(path string, info os.FileInfo, err error) error {

		fi2, err := os.Stat(path)
		switch {
		case err != nil:
			return nil
		case fi2.IsDir():
			return nil
		default:
			p := path[len(root):]
			found = append(found, p)
			return nil
		}
	}

	// Start walking!
	filepath.Walk(root, walker)
	return found, nil

}

// FilesInDirNotRecursive returns a list of all files in
// this root directory; but not in subdirectories.
// The root directory name is removed
// from the files; the returned names are relative to the root.
func FilesInDirNotRecursive(root string) ([]string, error) {
	found := []string{}

	fi, err := os.Stat(root)
	switch {
	case err != nil:
		return found, err
	case fi.IsDir():
		if strings.HasSuffix(root, "/") == false {
			root = root + "/"
		}
	default:
		return found, fmt.Errorf("%v: Not a directory", root)
	}

	// Create a callback function
	walker := func(path string, info os.FileInfo, err error) error {

		fi2, err := os.Stat(path)
		switch {
		case err != nil:
			return nil
		case fi2.IsDir():
			if path == root {
				return nil
			}
			return filepath.SkipDir
		default:
			p := path[len(root):]
			found = append(found, p)
			return nil
		}
	}

	// Start walking!
	filepath.Walk(root, walker)
	return found, nil

}
