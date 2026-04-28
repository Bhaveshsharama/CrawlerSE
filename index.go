package main

import (
	"encoding/gob"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

type Index struct {
	Inverted map[string][]string
	DF       map[string]int
	Total    int
	Pages    map[string]Page
	Version  int
}

type Pair struct {
	First  string
	Second int
}

const CurrentIndexVersion = 1

func BuildIndexFromScratch() (*Index, error) {

	wikiPages := Crawl(
		"https://en.wikipedia.org/wiki/Computer_science",
		1000,
		WikipediaPolicy{},
	)

	goPages := Crawl(
		"https://go.dev/doc/",
		1000,
		GenericPolicy{},
	)

	pages := append(wikiPages, goPages...)

	indexed := BuildIndex(pages)

	return indexed, nil
}

func BuildIndex(pages []Page) *Index {

	var idx Index

	idx.Pages = make(map[string]Page)

	for _, p := range pages {
		idx.Pages[p.URL] = p
	}

	idx.Inverted = make(map[string][]string) //val-> slice of urls
	idx.DF = make(map[string]int)

	for _, v := range pages {
		tokens := tokenise(v.Text)

		for _, w := range tokens {
			idx.Inverted[w] = append(idx.Inverted[w], v.URL)
		}
	}

	for key, val := range idx.Inverted {
		seen := make(map[string]bool)
		for i := 0; i < len(val); i++ {

			seen[val[i]] = true

		}
		idx.DF[key] = len(seen)
	}
	idx.Total = len(pages)
	idx.Version = CurrentIndexVersion

	return &idx

}

func (idx *Index) SaveIndex(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err

	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	if err := enc.Encode(idx); err != nil {
		return err

	}
	return nil

}

func LoadIndex(path string) (*Index, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var idx Index
	dec := gob.NewDecoder(f)
	if err := dec.Decode(&idx); err != nil {
		return nil, err
	}

	if idx.Version != CurrentIndexVersion {
		return nil, fmt.Errorf("index version mismatch")
	}

	return &idx, nil

}

func highlightTerm(text, term string) string {
	return strings.ReplaceAll(text, term, "["+term+"]")
}

func makeSnippet(text string, terms []string) string {
	original := text
	normalised := normalizeText(text)
	index := -1
	matchedTerm := ""

	for _, t := range terms {
		pos := strings.Index(normalised, t)
		if pos != -1 {
			index = pos
			matchedTerm = t
			break
		}
	}
	if index == -1 {
		if len(normalised) <= 120 {
			return normalised
		}
		return normalised[:120] + "..."
	}

	windowStart := index - 60
	if windowStart < 0 {
		windowStart = 0
	}

	windowEnd := index + len(matchedTerm) + 60
	if windowEnd > len(normalised) {
		windowEnd = len(normalised)
	}

	snippet := original[windowStart:windowEnd]

	if windowStart > 0 {
		snippet = "..." + snippet
	}

	if windowEnd < len(normalised) {
		snippet = snippet + "..."
	}
	for _, t := range terms {
		snippet = highlightTerm(snippet, t)
	}
	return snippet

}

func Search(idx *Index, query string, limit int, offset int) ([]Result, int) {

	scores := make(map[string]float64)

	seenterms := make(map[string]bool)
	results := []Result{}

	qNorm := normalizeText(query)
	tokens := tokenise(qNorm)

	totalLength := 0.0
	for _, page := range idx.Pages {
		totalLength += float64(len(strings.Fields(page.Text)))
	}
	avgDocLength := totalLength / float64(len(idx.Pages))

	logger := &Logger{}

	covered := make(map[string]map[string]bool)
	validTerms := []string{}

	for _, term := range tokens {

		_, ok := seenterms[term]
		if ok {
			continue
		}
		seenterms[term] = true

		urls, ok := idx.Inverted[term]

		if !ok {
			continue
		}

		df := idx.DF[term]

		if float64(df) >= math.Max(50, 0.8*float64(idx.Total)) {
			continue
		}

		validTerms = append(validTerms, term)

		tfmap := make(map[string]int)

		for _, u := range urls {
			tfmap[u]++
		}

		for key, val := range tfmap {
			title := normalizeText(idx.Pages[key].Title)

			idf := math.Log10((float64(idx.Total) + 1) / (float64(df) + 1))

			k1 := 1.2
			b := 0.75

			docLength := float64(len(strings.Fields(idx.Pages[key].Text)))

			numerator := float64(val) * (k1 + 1.0)
			denominator := float64(val) + k1*(1.0-b+b*(docLength/avgDocLength))

			bm25TF := numerator / denominator

			termScore := idf * bm25TF

			if strings.Contains(title, term) {
				termScore *= 2.5
			}
			scores[key] += termScore // title contribution

			if covered[key] == nil {
				covered[key] = make(map[string]bool)
			}
			covered[key][term] = true
		}
	}
	logger.DetectUnseenTerms(validTerms, idx.DF)

	totalTerms := len(validTerms)

	for doc := range scores {
		coveredCount := len(covered[doc])

		if totalTerms > 0 {
			coverage := float64(coveredCount) / float64(totalTerms)
			scores[doc] *= (1.0 + 0.5*coverage)
		}
	}

	for doc := range scores {
		combined := normalizeText(idx.Pages[doc].Title + " " + idx.Pages[doc].Text)

		if strings.Contains(combined, qNorm) {
			scores[doc] *= 1.5
		}
	}

	for key, val := range scores {
		results = append(results, Result{URL: key, Title: idx.Pages[key].Title, Count: val})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Count == results[j].Count {
			return results[i].URL < results[j].URL
		}
		return results[i].Count > results[j].Count
	})

	//pagination
	total := len(results)
	start := offset
	end := offset + limit
	if start > total {
		return []Result{}, total
	}

	if end > total {
		end = total
	}
	results = results[start:end]

	for i := range results {

		url := results[i].URL

		bestTerm := ""
		bestIDF := -1.0

		for term := range covered[url] {
			df := idx.DF[term]
			idf := math.Log10((float64(idx.Total) + 1) / (float64(df) + 1))

			if idf > bestIDF {
				bestIDF = idf
				bestTerm = term
			}
		}

		if bestTerm != "" {
			page := idx.Pages[url]
			titleNorm := normalizeText(page.Title)

			if strings.Contains(titleNorm, bestTerm) {
				snippet := page.Title
				for _, t := range tokens {
					snippet = highlightTerm(snippet, t)
				}
				results[i].Snippets = snippet
			} else {
				results[i].Snippets = makeSnippet(page.Text, tokens)
			}
		}
	}

	return results, total
}
