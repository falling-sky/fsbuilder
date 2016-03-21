package po

import "sync"

// Record is a single text translated
type Record struct {
	Comment string
	MsgID   string
	MsgStr  string
}

// MapStringRecord maps original strings to Records
type MapStringRecord map[string]*Record

// MapHeaders contains a list of headers from the "" element (first element).
type MapHeaders map[string]string

// File contains the map of strings for this translation.
type File struct {
	ByID              MapStringRecord
	InOrder           []string
	Headers           MapHeaders
	Language          string
	Locale            string
	Translated        int
	OutOf             int
	PercentTranslated string
	lock              sync.Mutex
}

// MapStringFile is a map of loaded translation files
type MapStringFile map[string]*File

// Files a collection of loaded translation files, plus the .pot file
type Files struct {
	Pot        *File
	ByLanguage MapStringFile
}
