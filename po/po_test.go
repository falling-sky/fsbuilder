package po

import (
	"testing"
)

func TestLoad(t *testing.T) {
	p, err := Load("../translations/dl/fr/falling-sky.fr_FR.po")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", p)

	var table = []struct {
		in  string
		out string
	}{
		{"slow", "lent"},
		{`Q: What do you mean by broken?`, `Q: Que voulez-vous dire par "disfonctionnel"?`},
	}

	for _, tt := range table {
		fetch, ok := p.ByID[tt.in]
		if !ok {
			t.Errorf("Error looking up: %v\n", tt.in)
			continue
		}
		if fetch.MsgStr != tt.out {
			t.Errorf("Bad comparison: checking %v, exppected %v, got %v\n", tt.in, tt.out, fetch.MsgStr)
			continue
		}
		t.Logf("success checking %v\n", tt.in)

	}
}

func TestLoadAll(t *testing.T) {
	multi, err := LoadAll("../translations/falling-sky.pot", "../translations/dl")
	if err != nil {
		t.Fatal(err)
	}
	//t.Logf("%#v", multi)
	for k, v := range multi.ByLanguage {
		t.Logf("%s: %v/%v", k, v.Translated, v.OutOf)
	}
	//t.Logf("%#v", multi.ByLanguage["pt_BR"])
}
