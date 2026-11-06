package main

import (
	"bufio"
	"strings"
	"testing"
)

func parse(t *testing.T, input string) [][]string {
	t.Helper()
	shares, err := parseShares(bufio.NewScanner(strings.NewReader(input)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return shares
}

func TestParseSharesSeparators(t *testing.T) {
	cases := map[string]string{
		"trailing blank line":         "# Share 0\na\nb\n\n# Share 1\nc\nd\n\n",
		"no trailing newline":         "a\nb\n\nc\nd",
		"multiple blank lines":        "\n\na\nb\n\n\n\nc\nd\n\n\n",
		"whitespace-only separator":   "a\nb\n  \t\nc\nd\n",
		"CRLF line endings":           "a\r\nb\r\n\r\nc\r\nd\r\n\r\n",
		"indented words and comments": "  # header\n  a \n\tb\n\n   # Share 1\nc\nd\n",
	}
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			shares := parse(t, input)
			if len(shares) != 2 {
				t.Fatalf("expected 2 shares, got %d: %q", len(shares), shares)
			}
			want := [][]string{{"a", "b"}, {"c", "d"}}
			for i := range want {
				if strings.Join(shares[i], " ") != strings.Join(want[i], " ") {
					t.Errorf("share %d: got %q, want %q", i, shares[i], want[i])
				}
			}
		})
	}
}

func TestParseSharesEmptyInput(t *testing.T) {
	for _, input := range []string{"", "\n\n", "# only a comment\n"} {
		_, err := parseShares(bufio.NewScanner(strings.NewReader(input)))
		if err == nil {
			t.Errorf("expected error for input %q", input)
		}
	}
}
