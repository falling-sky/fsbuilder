package po

import (
	"log"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

func Friendly(code string) string {
	l, e := language.Parse(code)
	if e != nil {
		log.Fatalf("Asked for friendly name for '%s', got error %v\n", code, e)
	}
	s := display.Self.Name(l)
	return s
}
