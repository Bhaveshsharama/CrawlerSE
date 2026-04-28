package main

import(
	"net/url"
	"strings"
)

type CrawlPolicy interface {
	Allows(u *url.URL) bool
	MaxPagesPerDomain() int
}

type WikipediaPolicy struct{} 

type GenericPolicy struct{}

func (p WikipediaPolicy) MaxPagesPerDomain() int {
    return 800
}

func (p GenericPolicy) MaxPagesPerDomain() int {
    return 250
}

func (p WikipediaPolicy) Allows(u *url.URL) bool {
	if u == nil {
		return false
	}
	host := u.Hostname()

	if (u.Scheme != "http" && u.Scheme != "https") || host == "" || u.Fragment != "" {  
		return false
	}

	if host!= "en.wikipedia.org" {
		return false
	}
	
	return isWikipediaArticle(u)
	
}

func isWikipediaArticle(u *url.URL) bool {
	path := u.Path

	if path == "/wiki/Main_Page" {
    	return false
	}

	if strings.Contains(path, ":") {
		return false
	}

	return len(path) > 6 && strings.HasPrefix(path, "/wiki/")
}

func (p GenericPolicy) Allows(u *url.URL) bool {
	if u == nil {
		return false
	}
	host := u.Hostname()

	if host!="go.dev" && host!="pkg.go.dev" {
		return false
	}

	if (u.Scheme != "http" && u.Scheme != "https") || host == "" || u.Fragment != "" { 
		return false
	}
	query:=u.RawQuery
	
	if strings.Contains(query,"utm_") || strings.Contains(query,"reply=") || strings.Contains(query,"share=") {
		return false
	}
	return true
}

