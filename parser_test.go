package main

import (
	"strings"
	"testing"
)

func TestParsePage_TitleAndLinks(t *testing.T) {
	html := `<html><head><title>My Page</title></head><body>
		<a href="/page1">L1</a>
		<a href="https://example.com/page2">L2</a>
	</body></html>`
	title, _, links := ParsePage(html, "https://example.com")

	if title != "My Page" {
		t.Errorf("title: got %q", title)
	}
	if len(links) < 2 {
		t.Fatalf("expected at least 2 links, got %d", len(links))
	}
	// relative link resolved
	found := false
	for _, l := range links {
		if strings.Contains(l, "example.com/page1") {
			found = true
		}
	}
	if !found {
		t.Error("relative link /page1 not resolved")
	}
}

func TestParsePage_FiltersInvisibleText(t *testing.T) {
	html := `<html><body><script>var x=1;</script><style>.a{}</style><p>visible</p></body></html>`
	_, text, _ := ParsePage(html, "https://example.com")
	norm := normalizeText(text)
	if !strings.Contains(norm, "visible") {
		t.Error("should contain visible text")
	}
	if strings.Contains(norm, "var") || strings.Contains(norm, "color") {
		t.Error("should not contain script/style content")
	}
}

func TestParsePage_WikipediaContentExtraction(t *testing.T) {
	html := `<html><body>
		<div id="nav">Nav junk</div>
		<div id="mw-content-text"><p>Article content here.</p></div>
		<div id="footer">Footer junk</div>
	</body></html>`
	_, text, _ := ParsePage(html, "https://en.wikipedia.org/wiki/Test")
	norm := normalizeText(text)
	if !strings.Contains(norm, "article content") {
		t.Error("should extract article content")
	}
	if strings.Contains(norm, "nav junk") || strings.Contains(norm, "footer junk") {
		t.Error("should exclude nav/footer when contentRoot found")
	}
}
