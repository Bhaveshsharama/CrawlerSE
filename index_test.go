package main

import (
	"strings"
	"testing"
)

func buildTestIndex() *Index {
	return BuildIndex([]Page{
		{URL: "https://example.com/go", Title: "Go Programming", Text: "go is a statically typed compiled programming language designed at google"},
		{URL: "https://example.com/python", Title: "Python Language", Text: "python is a high level general purpose programming language"},
		{URL: "https://example.com/rust", Title: "Rust Language", Text: "rust is a multi paradigm systems programming language focused on safety"},
	})
}

func TestBuildIndex(t *testing.T) {
	idx := buildTestIndex()
	if idx.Total != 3 {
		t.Errorf("Total: got %d, want 3", idx.Total)
	}
	if idx.DF["programming"] != 3 {
		t.Errorf("DF[programming]: got %d, want 3", idx.DF["programming"])
	}
	if idx.DF["golang"] > 0 && idx.DF["golang"] != 0 {
		t.Errorf("DF[golang] unexpected: %d", idx.DF["golang"])
	}
}

func TestSearch_RankingAndPagination(t *testing.T) {
	idx := buildTestIndex()

	results, total := Search(idx, "go programming", 10, 0)
	if total == 0 {
		t.Fatal("expected results for 'go programming'")
	}
	// results should be sorted descending
	for i := 1; i < len(results); i++ {
		if results[i].Count > results[i-1].Count {
			t.Error("results not sorted by score")
		}
	}
	// top result should be Go-related (title boost)
	if !strings.Contains(results[0].URL, "go") {
		t.Errorf("top result should be Go page, got %s", results[0].URL)
	}

	// pagination beyond results
	beyond, _ := Search(idx, "go", 10, 1000)
	if len(beyond) != 0 {
		t.Error("offset beyond total should return empty")
	}
}

func TestSearch_NoResults(t *testing.T) {
	idx := buildTestIndex()
	results, total := Search(idx, "xyznonexistent", 10, 0)
	if total != 0 || len(results) != 0 {
		t.Error("nonexistent query should return 0 results")
	}
}

func TestSaveAndLoadIndex(t *testing.T) {
	idx := BuildIndex([]Page{{URL: "https://a.com/1", Title: "Test", Text: "hello world"}})
	path := t.TempDir() + "/test.gob"

	if err := idx.SaveIndex(path); err != nil {
		t.Fatalf("SaveIndex: %v", err)
	}
	loaded, err := LoadIndex(path)
	if err != nil {
		t.Fatalf("LoadIndex: %v", err)
	}
	if loaded.Total != idx.Total || loaded.Version != idx.Version {
		t.Error("loaded index doesn't match saved index")
	}
}
