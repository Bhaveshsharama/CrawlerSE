# CrawlerSE

CrawlerSE is a highly concurrent, Go-based web crawler and full-text search engine. It crawls specified domains using configurable policies, builds an inverted index, and serves search results using a custom BM25 ranking algorithm via a minimalist, Google-style web interface. 

Built with standard library tools (zero framework dependencies on the backend, except for `golang.org/x/net/html` for parsing) to demonstrate core computer science and software engineering principles.

## Features

- **Concurrent Web Crawler**: Uses a worker pool pattern to fetch pages asynchronously while safely managing shared state (visited sets, domain frequency limits) with mutexes and atomics.
- **Configurable Crawl Policies**: Extensible interface (`CrawlPolicy`) to define domain-specific rules (e.g., distinguishing semantic Wikipedia pages from localized / talk namespaces).
- **Intelligent Parsing**: Extracts raw text while safely ignoring scripts, styles, and `<svg>` tags. Implements domain-specific extraction (e.g., focusing on `#mw-content-text` for Wikipedia to avoid indexing navbars/footers).
- **Text Processing pipeline**: Robust URL normalization, text sanitization, tokenization, and stopword filtering.
- **Inverted Index & BM25 Ranking**: Uses an inverted index for fast $O(1)$ term lookups, scoring documents via BM25 (handling term frequency saturation and document length normalization), enriched by Title Boosts and Query Coverage bonuses.
- **Snippets & Highlighting**: Auto-generates contextual snippets for search results focusing around query terms.
- **Minimalist Frontend**: Vanilla HTML/JS frontend featuring debounce, error states, and XSS protection.

## Architecture

1. **Crawler (`crawler.go`, `fetcher.go`)**: Seeds URLs into a channel. A set of goroutine workers concurrently fetch and parse pages until a generic or per-domain cap is reached.
2. **Parser (`parser.go`)**: Cleans HTML elements to find meaningful layout nodes using AST walking.
3. **Indexer & Ranker (`index.go`)**: Tokenizes the output, updates document frequencies (DF), stores inverted indexes, and saves the snapshot to disk as `index.gob` so subsequent runs can avoid re-crawling.
4. **API (`api.go`)**: Exposes an HTTP endpoint (`/search?q=...&limit=X&offset=Y`) that accepts queries, computes dynamic BM25 scores across the index, and paginates results.
5. **Frontend (`static/`)**: A clean UI to interface with the API.

## Getting Started

### Prerequisites
- Go 1.25+

### Installation & Running

1. Clone the repository:
   ```bash
   git clone https://github.com/Bhaveshsharama/CrawlerSE
   cd crawlerse
   ```

2. Download dependencies:
   ```bash
   go mod download
   ```

3. Run the engine:
   ```bash
   go run .
   ```
   *Note: On the first run, the engine will crawl the web (seeded at Wikipedia and Go docs) and build `index.gob`. This may take a moment. Subsequent runs will be significantly faster by loading the cached index.*

4. Open your browser and navigate to:
   ```
   http://localhost:8080
   ```

## Running Tests
The project features a comprehensive test suite covering the parser, indexing, utilities, and BM25 relevance edge-cases.

```bash
go test -v ./...
```

## Technical Decisions & Future Improvements

- **Why Go?** Native concurrency primitives (goroutines and channels) make writing robust crawler worker pools straightforward and highly performant. 
- **Why BM25 instead of TF-IDF?** BM25 introduces term-frequency saturation, meaning that mentioning a keyword 10 times is better than once, but not 10 times better. It balances document lengths so overwhelmingly sparse or dense files don't hijack the scoring.
- **Future Improvements**:
    - **`robots.txt` compliance**: Dynamically fetching and parsing `robots.txt` using the `temoto/robotstxt` package before visiting a domain.
    - **Persistence**: Swapping the in-memory `map[string][]string` / `gob` file for a real disk-backed distributed store like ElasticSearch or Redis if scaling beyond a few thousand pages.
    - **Distributed Workers**: Implementing a message queue (RabbitMQ / Kafka) to decouple crawling nodes from the indexing node.