package main

import (
	"net/url"
	"testing"
)

func TestWikipediaPolicy(t *testing.T) {
	p := WikipediaPolicy{}
	allow := []string{"https://en.wikipedia.org/wiki/Computer_science"}
	reject := []string{
		"https://en.wikipedia.org/wiki/Main_Page",
		"https://en.wikipedia.org/wiki/Talk:Go",
		"https://en.wikipedia.org/wiki/Category:Languages",
		"https://example.com/wiki/Test",
	}
	for _, raw := range allow {
		u, _ := url.Parse(raw)
		if !p.Allows(u) {
			t.Errorf("should allow %s", raw)
		}
	}
	for _, raw := range reject {
		u, _ := url.Parse(raw)
		if p.Allows(u) {
			t.Errorf("should reject %s", raw)
		}
	}
	if p.Allows(nil) {
		t.Error("should reject nil")
	}
}

func TestGenericPolicy(t *testing.T) {
	p := GenericPolicy{}
	allow := []string{"https://go.dev/doc/", "https://pkg.go.dev/fmt"}
	reject := []string{"https://google.com/search", "https://go.dev/blog?utm_source=tw"}
	for _, raw := range allow {
		u, _ := url.Parse(raw)
		if !p.Allows(u) {
			t.Errorf("should allow %s", raw)
		}
	}
	for _, raw := range reject {
		u, _ := url.Parse(raw)
		if p.Allows(u) {
			t.Errorf("should reject %s", raw)
		}
	}
}
