package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
)

// Record contains configuration options
type Record struct {
	Directories struct {
		TemplateDir    string
		ImagesDir      string
		TransparentDir string
		PoDir          string
		OutputDir      string
	}
	Processors struct {
		Note   []string
		JS     []string
		CSS    []string
		HTML   []string
		PHP    []string
		Apache []string
	}
	Map     map[string]string
	Options struct {
		MaxThreads int
	}
}

// Defaults will update a config record with safe defaults for any missing values
func (r *Record) Defaults() {
	if r.Directories.TemplateDir == "" {
		r.Directories.TemplateDir = "templates"
	}
	if r.Directories.ImagesDir == "" {
		r.Directories.ImagesDir = "images"
	}
	if r.Directories.TransparentDir == "" {
		r.Directories.TransparentDir = "transparent"
	}
	if r.Directories.PoDir == "" {
		r.Directories.PoDir = "translations"
	}
	if r.Directories.OutputDir == "" {
		r.Directories.OutputDir = "output"
	}

	if len(r.Processors.Note) == 0 {
		r.Processors.Note = []string{
			"Macros available:",
			"[NAME] will simply be index.html.en_US, index.js.en_US, or comment.php",
			"[NAMEGZ] will simply be index.html.gz.en_US, index.js.gz.en_US, or comment.php.gz",
			"For convenience:",
			"[INPUT] will be identical to [NAME].orig - and is written to disk at the start of the commands.",
			"[OUTPUT] will be idetnical to [NAME]",
			"All processors must be defined.  At minimum simply use mv [NAME].orig NAME",
		}
	}

	if len(r.Processors.JS) == 0 {
		r.Processors.JS = []string{
		//			`uglifyjs2  [NAME].orig -o [NAME] -c --warnings=false   --source-map [NAME].map   --stats`,
		//			`mv [NAME].orig [NAME]`,
		//			`gzip -f -9 -Sgz  < [NAME]  > [NAMEGZ]`,
		}
	}
	if len(r.Processors.CSS) == 0 {
		r.Processors.CSS = []string{
		//			`cssmin < [NAME].orig > [NAME]`,
		//			`mv [NAME].orig [NAME]`,
		//			`gzip -f -9 -Sgz  < [NAME]  > [NAMEGZ]`,
		}
	}
	if len(r.Processors.HTML) == 0 {
		r.Processors.HTML = []string{
		//		`tidy -quiet -indent -asxhtml -utf8 -w 120 --show-warnings false < [NAME].orig > [NAME]`,
		//		`sed < [NAME] 's#/index.js#/index.js.gz#' | gzip -f -9 -Sgz  > [NAMEGZ]`,
		}
	}
	if len(r.Processors.PHP) == 0 {
		r.Processors.PHP = []string{
		//			`mv [NAME].orig [NAME]`,
		}
	}
	if len(r.Processors.Apache) == 0 {
		r.Processors.Apache = []string{
		//			`mv [NAME].orig [NAME]`,
		}
	}

	if r.Map == nil {
		r.Map = make(map[string]string)
	}
	if len(r.Map) == 0 {
		r.Map["dot.htaccess"] = ".htaccess"
		r.Map["ip.htaccess"] = "ip/.htaccess"
		r.Map["images.htaccess"] = "images/.htaccess"
		r.Map["images-nc.htaccess"] = "images-nc/.htaccess"
	}

}

// Load a config file, return it after adjusting for defaults
func Load(filename string) (*Record, error) {
	r := &Record{}

	// If a filename is specified, load it.
	if filename != "" {
		b, e := ioutil.ReadFile(filename)
		if e != nil {
			return r, e
		}
		e = json.Unmarshal(b, r)
		if e != nil {
			return r, e
		}
	}
	r.Defaults()
	return r, nil
}

func (r *Record) String() string {
	b, e := json.MarshalIndent(r, "", "\t")
	if e != nil {
		log.Fatal(e)
	}

	b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
	b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
	b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)

	return string(b)
}

// Return a sample config with defaults
func Example() string {
	r := &Record{}
	r.Defaults()
	return r.String()
}
