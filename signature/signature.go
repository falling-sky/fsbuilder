package signature

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"path/filepath"

	"github.com/falling-sky/builder/fileutil"
)

func ScanDir(directory string, otherstuff ...string) string {
	h := md5.New()

	log.Printf("ScanDir(%s)", directory)
	files, err := fileutil.FilesInDirRecursive(directory)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		e := filepath.Ext(file)
		if e != ".html" && e != ".js" && e != ".htaccess" && e != ".inc" && e != ".example" && e != ".php" {
			continue
		}

		fn := directory + "/" + file
		//	log.Printf("scanning %s\n", fn)

		// This will cache, saving a trip for other jobs
		content, err := fileutil.ReadFile(fn)
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(h, content)

	}

	// Any other stuff we passed, to bias a signature.
	// Such as languages.
	for _, s := range otherstuff {
		io.WriteString(h, s)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
