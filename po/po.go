package po

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var example = `
#: faq_whyipv6.html
msgid "Note this is in addition to any NAT you do at home."
msgstr ""

#: faq_whyipv6.html
msgid "Q: So, why worry? NAT will work, right? I use NAT at home today after all.."
msgstr ""
`

var reWHITESPACE = regexp.MustCompile(`\s+`)

func unquote(s string) (string, error) {
	if len(s) == 0 {
		return s, nil
	}
	s = strings.TrimSpace(s)
	return strconv.Unquote(s)
}

// Languages returns the list of locales loaded in the combined *Files object
func (combined *Files) Languages() []string {
	ret := []string{}

	for k := range combined.ByLanguage {
		ret = append(ret, k)
	}
	sort.Strings(ret)
	return ret
}

// GetLocale simply returns the locale name; ie en_US or pt_BR
func (f *File) GetLocale() string {
	s := f.Locale
	return s
}

// GetLang returns the lowercase name string; ie en or pt
func (f *File) GetLang() string {
	s := f.Locale
	p := strings.Split(s, "_")
	return p[0]
}

// GetLangUC  returns the uppercase name string; ie EN or PT
func (f *File) GetLangUC() string {
	s := f.Locale
	p := strings.Split(s, "_")
	return strings.ToUpper(p[0])
}

// GetLang returns the lowercase name string; ie en or pt
func (f *File) GetLangName() string {
	s := f.Language
	s = Friendly(s)
	return s
}

// GetLangPercentTranslated returns what percentage of the translation is done
func (f *File) GetLangPercentTranslated() string {
	s := f.PercentTranslated
	return s
}

// Translate takes a given input text, and returns back
// either the translated text, or the original text again.
func (f *File) Translate(input string, escapequotes bool) string {

	// Canonicalize.
	// Remove redundant, leading, and trailing whitespace.
	input = strings.TrimSpace(input)
	input = reWHITESPACE.ReplaceAllString(input, " ")

	if input == "lang" {
		return f.GetLang()
	}
	if input == "langUC" {
		return f.GetLangUC()
	}
	if input == "locale" {
		return f.GetLocale()
	}
	if input == "langname" {
		return f.GetLangName()
	}
	if input == "percenttranslated" {
		return f.GetLangPercentTranslated()
	}

	newtext := input

	if found, ok := f.ByID[input]; ok {
		c := found.MsgStr
		if c != "" {
			newtext = c
		}
	}

	if escapequotes {
		newtext = strings.Replace(newtext, `"`, `\"`, -1)
		newtext = strings.Replace(newtext, `'`, `\'`, -1)
	}
	// TODO escapequotes
	// Perl does this:
	//         $text =~ s/(?<![\\])"/\\"/g;
	//    $text =~ s/(?<![\\])'/\\'/g;
	// GO does not do look-behind assertions
	return newtext
}

// Translate takes a given input text, and returns back
// either the translated text, or the original text again.
func (f *File) Add(input string, context string, escapequotes bool) {

	//	log.Printf("po file Add(%s)\n", input)
	// Canonicalize.
	// Remove redundant, leading, and trailing whitespace.
	input = strings.TrimSpace(input)
	input = reWHITESPACE.ReplaceAllString(input, " ")

	//	log.Printf("po Add input=%s context=%s escape=%v\n", input, context, escapequotes)

	// Skip these, these will be dynamically responded to.
	if input == "lang" || input == "langUC" || input == "locale" {
		return
	}

	f.lock.Lock()
	defer f.lock.Unlock()

	_ = "breakpoint"

	if v, ok := f.ByID[input]; ok == false {
		// Not yet set?  Let's do so.
		//	fmt.Printf("DEBUG Saving Comment=%s MsgID=%s\n", context, input)
		f.ByID[input] = &Record{
			Comment: context,
			MsgID:   input,
			MsgStr:  "",
		}
		f.InOrder = append(f.InOrder, input)
	} else if strings.Contains(v.Comment, "not-used") {
		v.Comment = context // Update context
	} else {
		//	fmt.Printf("DEBUG ALREADY HAVE %#v\n", f.ByID[input])

	}

	/*
		if escapequotes {
			newtext = strings.Replace(newtext, `"`, `\"`, -1)
			newtext = strings.Replace(newtext, `'`, `\'`, -1)
		}
	*/

}

// ApacheAddLanguage  Generates the Apache "AddLanguage" text
func (f *Files) ApacheAddLanguage() string {
	list := append([]string{"en_US"}, f.Languages()...)
	text := ""
	seen := make(map[string]bool)

	add := func(s string) {
		if seen[s] == false {
			text = text + s + "\n"
			seen[s] = true
		}
	}

	for _, locale := range list {
		parts := strings.Split(locale, "_")
		if len(parts) > 0 {
			add(fmt.Sprintf("AddLanguage %s .%s", parts[0], locale))
		}
		dashed := strings.Replace(locale, "_", "-", -1)
		add(fmt.Sprintf("AddLanguage %s .%s", dashed, locale))

	}

	return text

}
