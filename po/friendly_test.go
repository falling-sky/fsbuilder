package po

import "testing"

func TestFriendly(t *testing.T) {
	var table = []struct {
		in  string
		out string
	}{
		{"af_ZA", "English"},
		{"az_AZ", "English"},
		{"cs_CZ", "English"},
		{"de_DE", "English"},
		{"el_GR", "English"},
		{"es_ES", "English"},
		{"fi_FI", "English"},
		{"fr_FR", "English"},
		{"hr_HR", "English"},
		{"hu_HU", "English"},
		{"it_IT", "English"},
		{"ja_JP", "English"},
		{"nb_NO", "English"},
		{"nl_NL", "English"},
		{"pl_PL", "English"},
		{"pt_BR", "English"},
		{"ro_RO", "English"},
		{"ru_RU", "English"},
		{"sk_SK", "English"},
		{"sq_AL", "English"},
		{"sv_SE", "English"},
		{"tr_TR", "English"},
		{"zh_CN", "English"},
		{"zh_TW", "English"},
	}
	for _, tt := range table {
		found := Friendly(tt.in)
		if found == tt.out {
			t.Logf("%s=%s (ok)\n", tt.in, found)
		} else {
			t.Errorf("%s=%s (bad)\n", tt.in, found)
		}
	}
}
