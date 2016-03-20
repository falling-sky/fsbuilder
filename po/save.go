package po

import (
	"bytes"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func PoQuote(w *bytes.Buffer, label string, content string) {
	w.WriteString(label)
	w.WriteString(" ")

	if strings.Contains(content, "\n") {
		w.WriteString(strconv.Quote(""))
		w.WriteString("\n")
		lines := strings.SplitAfter(content, "\n")
		for _, line := range lines {
			if len(line) > 0 {
				w.WriteString(strconv.Quote(line))
				w.WriteString("\n")
			}
		}
	} else {
		w.WriteString(strconv.Quote(content))
		w.WriteString("\n")
	}
}

// Load a .PO file into memory.
func (f *File) Save(fn string) error {
	log.Printf("Generating %s\n", fn)
	f.ByID[""] = &Record{
		MsgID: "",
		MsgStr: `Project-Id-Version: PACKAGE VERSION
PO-Revision-Date: YEAR-MO-DA HO:MI +ZONE
Last-Translator: Unspecified Translator <jfesler+unspecified-translator@test-ipv6.com>
Language-Team: LANGUAGE <v6code@test-ipv6.com>
MIME-Version: 1.0
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: 8bit
`,
	}

	// Start new output buffer
	b := &bytes.Buffer{}

	// Prepend "" into the order
	f.InOrder = append([]string{""}, f.InOrder...)

	for _, str := range f.InOrder {
		r := f.ByID[str]
		if r.Comment != "" {
			PoQuote(b, "#:", r.Comment)
		}
		PoQuote(b, "msgid", r.MsgID)
		PoQuote(b, "msgstr", r.MsgStr)
		b.WriteString("\n")
	}
	err := ioutil.WriteFile(fn, b.Bytes(), 0644)
	return err

}
