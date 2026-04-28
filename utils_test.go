package main

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct{ input, want string }{
		{"https://example.com/page#section", "https://example.com/page"},
		{"https://example.com/path/", "https://example.com/path"},
		{"https://EN.WIKIPEDIA.ORG/wiki/Test", "https://en.wikipedia.org/wiki/Test"},
		{"", ""},
	}
	for _, tc := range tests {
		if got := normalizeURL(tc.input); got != tc.want {
			t.Errorf("normalizeURL(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestNormalizeText(t *testing.T) {
	tests := []struct{ input, want string }{
		{"  Hello,  WORLD!!!  ", "hello world"},
		{"multiple   spaces\tand\nnewlines", "multiple spaces and newlines"},
		{"", ""},
	}
	for _, tc := range tests {
		if got := normalizeText(tc.input); got != tc.want {
			t.Errorf("normalizeText(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestTokenise(t *testing.T) {
	got := tokenise("the quick brown fox is in the forest")
	expected := map[string]bool{"quick": true, "brown": true, "fox": true, "forest": true}
	if len(got) != len(expected) {
		t.Fatalf("tokenise: got %v, want keys of %v", got, expected)
	}
	for _, tok := range got {
		if !expected[tok] {
			t.Errorf("unexpected token %q", tok)
		}
	}
}
