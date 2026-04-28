package main

import (
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"
)

func Crawl(startURL string, maxPages int,policy CrawlPolicy) []Page {

	var pages []Page
	
	const workerCount = 3

	var activeWork int64

	urlChan := make(chan string, 1000)
	done:=make(chan struct{})
	var shutdownonce sync.Once

	visited := make(map[string]bool)
	seenFrontier:=make(map[string]bool)
	domainCount := make(map[string]int)

	var mu sync.Mutex

	var workerWG sync.WaitGroup

	shutdown:= func(){
		shutdownonce.Do(func() {
			close(done)
			
		})
	}
	
	worker := func() {
		defer workerWG.Done()

		for {
			select {
			case <-done:
				return
			case rawURL := <- urlChan:
				atomic.AddInt64(&activeWork, 1)
				func() {
				defer func() {
    				if atomic.AddInt64(&activeWork, -1) == 0 {
        				shutdown()
    				}
				}()
				
			

			norm := normalizeURL(rawURL)
			if norm == "" {
				
				return
			}

			u, err := url.Parse(norm)
			if err != nil {    
				
				return
			}
			if !policy.Allows(u) {				
				return
			}

			domain := u.Hostname()

			mu.Lock()

			if visited[norm] || len(pages) >= maxPages || domainCount[domain] >= policy.MaxPagesPerDomain() {
    			mu.Unlock()
				return
			}

			
			visited[norm]=true
			mu.Unlock()

			htmlStr, err := Fetch(norm)
			if err != nil {
				fmt.Println(err)
				return
			}
			

			title, text, links := ParsePage(htmlStr, norm)

			mu.Lock()
			pages = append(pages, Page{
				URL:   norm,
				Title: title,
				Text:  normalizeText(text),
				 
			})
			domainCount[domain]++
			
			reachedLimit := len(pages) >= maxPages
			mu.Unlock()
			if reachedLimit {
				shutdown()
				
			}

			for _, link := range links {
				l := normalizeURL(link)
				if l == "" {
					continue
				}

				lu, err := url.Parse(l)
				if err!=nil {
					continue
				}
				mu.Lock();
				if seenFrontier[l] {
					mu.Unlock()
					continue
				}
				seenFrontier[l]=true        
				mu.Unlock()

				if !policy.Allows(lu) {
					
					continue
				}
				select {
				case <-done:
					return
				default:
				}

				select {
				case urlChan <- l:     
				default:
				
				}
			}
		}()
		}
		}
	}

	// start workers
	for i := 0; i < workerCount; i++ {
		workerWG.Add(1)
		go worker()
	}

	// seed
	
	urlChan <- startURL

	workerWG.Wait()

	return pages
}