package repository

import "testing"

func TestLikePattern(t *testing.T) {
	cases := map[string]string{
		"":              "%%",             // match anything (empty never reaches search in practice)
		"docker":        "%docker%",       // single word
		"docker port":   "%docker%port%",  // words joined in order
		"  go   test  ": "%go%test%",      // extra whitespace collapsed
		"git_log":       `%git\_log%`,     // underscore escaped (not a wildcard)
		"50%":           `%50\%%`,         // percent escaped
		`a\b`:           `%a\\b%`,         // backslash escaped
	}
	for in, want := range cases {
		if got := likePattern(in); got != want {
			t.Errorf("likePattern(%q) = %q, want %q", in, got, want)
		}
	}
}
