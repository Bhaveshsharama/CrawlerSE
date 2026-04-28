package main

type Page struct {
	URL   string
	Title string
	Text  string
}

type Result struct {
	URL      string  `json:"url"`
	Title    string  `json:"title"`
	Count    float64 `json:"score"`
	Snippets string  `json:"snippets"`
}

type SearchResponse struct {
	Total   int      `json:"total"`
	Results []Result `json:"results"`
	TimeMS  int64    `json:"time_ms"`
}
