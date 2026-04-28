package main

import (
	"net/url"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

func normalizeURL(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}

	parsed.Fragment = ""
	parsed.Host = strings.ToLower(parsed.Host)

	parsed.Path = strings.TrimSuffix(parsed.Path, "/")

	return parsed.String()
}

func normalizeText(s string) string {
	s = strings.ToLower(s)

	var result strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			result.WriteRune(r)
		}
	}

	return strings.Join(strings.Fields(result.String()), " ")
}

func isVisibleText(n *html.Node) bool {
	if n.Parent == nil {
		return false
	}

	switch n.Parent.Data {
	case "script", "style", "noscript", "svg", "meta", "head":
		return false
	}
	return true
}

var stopwords = map[string]bool{
	"the": true, "is": true, "and": true, "of": true,
	"to": true, "in": true, "for": true, "on": true,
	"with": true, "as": true, "by": true, "an": true,
}

func tokenise(text string) []string {
	normalized := normalizeText(text)
	tokens := strings.Fields(normalized)

	finaltokens := make([]string, 0, len(tokens))
	for _, t := range tokens {
		if stopwords[t] == true {
			continue
		}
		finaltokens = append(finaltokens, t)
	}
	return finaltokens
}
