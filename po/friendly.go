package po

import (
	"golang.org/x/text/language"
	"golang.org/x/text/language/display"
)

func Friendly(code string) string {
	l := language.MustParse(code)
	s := display.Self.Name(l)
	return s
}
