package job

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"text/template"

	"github.com/falling-sky/fsbuilder/config"
	"github.com/falling-sky/fsbuilder/fileutil"
	"github.com/falling-sky/fsbuilder/gitinfo"
	"github.com/falling-sky/fsbuilder/po"
)

// rePROCESS matches on   [% PROCESS "filename" %]
// and captures the entire template directivel as well as the inside filename.
var rePROCESS = regexp.MustCompile(`\[\%\s*PROCESS\s*"(.*?)"\s*\%\]`)
var reTRANSLATE = regexp.MustCompile(`(?ms){{(.*?)}}`)

// PostType describes a directory, and how to process it.
type PostInfoType struct {
	Directory   string
	Extension   string
	PostProcess []string
	EscapeQuote bool
	MultiLocale bool
	Compress    bool
}

// QueueItem represents a single job to be queued, and ran as capacity allows.
// This is so we can generate the work list up front; and then pace out the work
// based on number of avaialble CPUs.
type QueueItem struct {
	Config   *config.Record
	RootDir  string
	Filename string
	PoFile   *po.File
	PotFile  *po.File
	Data     *TemplateData
	PostInfo PostInfoType
}

// QueueTracker is an object for managing QueueItem jobs.
type QueueTracker struct {
	Channel chan *QueueItem
	WG      *sync.WaitGroup
}

// TemplateData is passed when adding the job to the queue.
// This is used by Go's text/template to extract info before expansion.
// Because templates are cached after being processed,
// and before translation, we should not place anything here that is
// locale specific.
type TemplateData struct {
	GitInfo      *gitinfo.GitInfo
	PoMap        po.MapStringFile
	Basename     string
	AddLanguage  string
	DirSignature string
}

// ParsedCacheType provides properly mutex locked cache access to
// the expanded (but untranslated) templates.
type ParsedCacheType struct {
	lock   sync.RWMutex
	byname map[string]string
}

// ParsedCache holds the actual cache of expanded (but not translated) templates.
var ParsedCache ParsedCacheType

func init() {
	ParsedCache.byname = make(map[string]string)
}

// GrabContent grabs a file.  Takes into account the QueueItem variables
// such as the iput directory path.  The file is cached for future requests.
func GrabContent(qi *QueueItem) string {
	topName := qi.RootDir + "/" + qi.Filename
	// log.Printf("GrabContent(%s)  (%s)\n", qi.Filename, qi.PoFile.Language)

	grab := func(fn string) string {
		//		log.Printf("GrabContent(%s)  (%s) (fn=%s)\n", qi.Filename, qi.PoFile.Language, fn)

		fullname := qi.RootDir + "/" + fn
		c, err := fileutil.ReadFile(fullname)
		if err != nil {
			log.Fatalf("tried to load %s (via %s): %s", fullname, topName, err)
		}
		//		log.Printf("read %v (%v bytes)\n", fullname, len(c))

		UpdatePot(qi, c, fn)

		return c
	}

	content := grab(qi.Filename)
	// Do we see PROCESS lines?
	for {
		matches := rePROCESS.FindStringSubmatch(content)
		if len(matches) == 0 {
			break
		}
		if len(matches) < 2 {
			log.Fatalf("I don't know what happened, but %s is interesting", matches[0])
		}
		wrapperString := matches[0]
		insideName := matches[1]
		newContent := grab(insideName)
		content = strings.Replace(content, wrapperString, newContent, -1)
	}
	return content
}

// ProcessTemplate runs text.Template against the given text.
// Note we use [% %]  for text.Template directorives, since these
// are fewer than translations. And we prefer to do translations
// without the template ugliness.
func ProcessTemplate(qi *QueueItem, content string) string {
	topName := qi.RootDir + "/" + qi.Filename

	// Do we need any custom functions?
	FuncMap := make(template.FuncMap)
	FuncMap["EXAMPLE"] = func(name string) (string, error) {
		//log.Printf("PROCESS: %v\n", name)
		return "", nil
	}

	// Parse the template.  Just looks for markers and implied commands.
	root := template.New(qi.Filename).Delims(`[%`, `%]`).Funcs(FuncMap)
	tmpl, err := root.Parse(content)
	if err != nil {
		log.Fatalf("Parsing template for %v: %v", topName, err)
	}

	// Execute the template.
	wr := &bytes.Buffer{}
	err = tmpl.Execute(wr, qi.Data)
	if err != nil {
		log.Fatalf("Executing template for %v: %v", topName, err)
	}

	return string(wr.Bytes())
}

func UpdatePot(qi *QueueItem, content string, fn string) {
	//	log.Printf("UpdatePot fn=%s\n", fn)
	for {

		matches := reTRANSLATE.FindStringSubmatch(content)
		if len(matches) == 0 {
			break
		}
		if len(matches) < 2 {
			log.Fatalf("I don't know what happened, but %s is interesting", matches[0])
		}
		wrapperString := matches[0]
		insideName := matches[1]

		//		log.Printf("UpdatePot inside=%s fn=%s escape=%v\n", insideName, fn, qi.EscapeQuotes)

		qi.PotFile.Add(insideName, fn, qi.PostInfo.EscapeQuote)

		content = strings.Replace(content, wrapperString, insideName, -1)

	}
}

// TranslateContent  looks for {{ text }} and replaces it with
// either translated text, or the original text.
func TranslateContent(qi *QueueItem, content string) string {
	for {

		matches := reTRANSLATE.FindStringSubmatch(content)
		if len(matches) == 0 {
			break
		}
		if len(matches) < 2 {
			log.Fatalf("I don't know what happened, but %s is interesting", matches[0])
		}
		wrapperString := matches[0]
		insideName := matches[1]

		//	log.Printf("grabbing %v\n", insideName)
		newContent := qi.PoFile.Translate(insideName, qi.PostInfo.EscapeQuote)

		//	log.Printf("Replacing %s with %s\n", wrapperString, newContent)

		content = strings.Replace(content, wrapperString, newContent, -1)

	}
	return content
}

func ProcessContentFancy(qi *QueueItem, content string) {

	tasks := qi.PostInfo.PostProcess

	basename := qi.Filename
	if t, ok := qi.Config.Map[qi.Filename]; ok {
		basename = t
	}

	// Prepare the macros that we support for running external commands.
	macros := make(map[string]string)
	macros["NAME"] = basename
	macros["NAMEGZ"] = macros["NAME"] + ".gz"
	if qi.PostInfo.MultiLocale == true {
		macros["NAME"] = macros["NAME"] + "." + qi.PoFile.Locale
		macros["NAMEGZ"] = macros["NAMEGZ"] + "." + qi.PoFile.Locale
	}
	macros["INPUT"] = macros["NAME"] + ".orig"
	macros["OUTPUT"] = macros["NAME"]

	// Supported varables are:
	// INPUT                  - same as [FILENAME].orig - the file written by builder before calling 3rd party fancy tools
	// FILENAME or OUTPUT     - Final output filename (with path)
	// FILENAMEGZ or OUTPUTGZ - Final output filename, gzipped (with path)
	// BASENAME               - Final output name (no path). Use this if inserting text into files.
	// BASENAMEGZ             - Final output name, gzipped (no path).  Use this if inserting text into files.

	fixup := func(s string) string {
		for k, v := range macros {
			kk := "[" + k + "]"
			s = strings.Replace(s, kk, v, -1)
		}
		return s
	}

	// First, write the file to disk.
	outputfilename := qi.Config.Directories.OutputDir + "/" + macros["INPUT"]
	os.MkdirAll(filepath.Dir(outputfilename), 0755)

	err := ioutil.WriteFile(outputfilename, []byte(content), 0755)
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("wrote %s etc (%v bytes)\n", outputfilename, len(content))

	// Post processing defined from the config file, 3rd party tools
	for _, task := range tasks {
		runcmd := fixup(task)
		//log.Printf("running: %s\n", runcmd)

		shellscript := bytes.NewBufferString(runcmd)
		stderr := &bytes.Buffer{}

		c := exec.Cmd{}
		c.Path = "/bin/sh"
		c.Dir = qi.Config.Directories.OutputDir
		c.Stdin = shellscript
		c.Stderr = stderr
		// log.Printf("About to run: %#v\n", runcmd)
		e := c.Run()

		// HACK HACK HACK ignore tidy exit code 1
		if e != nil {
			if strings.HasPrefix(runcmd, "tidy ") {
				if e.Error() == "exit status 1" {
					e = nil
				}
			}
		}
		if e != nil {
			log.Printf("stderr: %s\n", stderr.String())
			log.Fatalf("While running %#v .. got: %#v\n", runcmd, e.Error())
		}
	}

}

func ProcessContent(qi *QueueItem, content string) {

	// See if there are commands specified. IF so, run those.
	tasks := qi.PostInfo.PostProcess
	if len(tasks) > 0 {
		ProcessContentFancy(qi, content)
		return
	}

	basename := qi.Filename
	if t, ok := qi.Config.Map[qi.Filename]; ok {
		basename = t
	}

	// Otherwise, do writes directly, and do our own compression.
	uncompressed := qi.Config.Directories.OutputDir + "/" + basename
	compressed := qi.Config.Directories.OutputDir + "/" + basename + ".gz"
	if qi.PostInfo.MultiLocale == true {
		uncompressed = uncompressed + "." + qi.PoFile.Locale
		compressed = compressed + "." + qi.PoFile.Locale
	}

	// Make sure the directory exists.
	// We may need to keep track of this;
	// do we really want to do this 1000+ times?
	os.MkdirAll(filepath.Dir(uncompressed), 0755)

	err := ioutil.WriteFile(uncompressed, []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("wrote %s etc (%v bytes)\n", outputfilename, len(content))

	if qi.PostInfo.Compress {
		if strings.HasSuffix(qi.Filename, ".html") {
			content = strings.Replace(content, `src="/index.js`, `src="/index.js.gz`, -1)
			content = strings.Replace(content, `href="/index.css`, `href="/index.css.gz`, -1)
		}

		// Compress in memory
		b := &bytes.Buffer{}
		w, err := gzip.NewWriterLevel(b, gzip.BestCompression)
		if err != nil {
			log.Fatal(err)
		}
		w.Write([]byte(content))
		w.Close()

		// And write
		err = ioutil.WriteFile(compressed, b.Bytes(), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

}

// RunJob takes a single QueueItem, and expands, translates, optimizes,
// and writes files for that single file for a single language.  These are spoon-fed
// by RunQueue.
func RunJob(qi *QueueItem) {
	// log.Printf("RunJob Filename=%s PoLang=%s\n", qi.Filename, qi.PoFile.Language)
	readFilename := qi.RootDir + "/" + qi.Filename

	var content string
	ParsedCache.lock.Lock()
	if c, ok := ParsedCache.byname[readFilename]; ok {
		// log.Printf("cached: %s", readFilename)
		content = c
	} else {
		//log.Printf("not cached: %s", readFilename)
		content = GrabContent(qi)
		content = ProcessTemplate(qi, content)
		ParsedCache.byname[readFilename] = content
	}
	ParsedCache.lock.Unlock()

	// TODO process translations
	content = TranslateContent(qi, content)
	ProcessContent(qi, content)

}

// RunQueue is a goroutine that listens to a channel for jobs.
// If jobs are accepted, they are given to RunJob.
func (qt *QueueTracker) RunQueue() {
	for {
		job, ok := <-qt.Channel
		if ok {
			RunJob(job)  // Run the job.
			qt.WG.Done() // Decrement WaitGroup counter
		} else {
			return
		}
	}
}

// Add a job to the queue.  Sends it to the channel.
func (qt *QueueTracker) Add(qi *QueueItem) {
	qt.WG.Add(1)     // Increment the WaitGroup counter.
	qt.Channel <- qi // Put the job in the queue.
}

// Wait will wait for all existing jobs to finish.
func (qt *QueueTracker) Wait() {
	log.Printf("Waiting for queued jobs to finish\n")
	qt.WG.Wait()
}

// StartQueue will start a goroutine for jobs, and return
// a handle to be used for adding and waiting on jobs.
func StartQueue(maxjobs int) *QueueTracker {
	qt := &QueueTracker{}
	qt.Channel = make(chan *QueueItem, 10000)
	qt.WG = &sync.WaitGroup{}

	if maxjobs == 0 {
		maxjobs = runtime.NumCPU()
	}
	for i := 0; i < maxjobs; i++ {
		go qt.RunQueue()
	}

	return qt
}
