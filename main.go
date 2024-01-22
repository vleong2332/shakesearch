package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handleSearch(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("shakesearch available at http://localhost:%s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		results, searchError := searcher.Search(query[0])
		if searchError != nil {
			return
		}
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		encodingError := enc.Encode(results)
		if encodingError != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New(dat)
	return nil
}

/**
 * Returns (results, hasMore, error)
 */
func (s *Searcher) Search(query string) ([]string, error) {
	/**
	 * TODO: Current trade-off is case-insensitivity vs. correctness. Something's wrong with punctuation:
	 * "posting is no need. " works; "posting is no need. O" doesn't.
	 * Tried regexp.QuoteMeta(query) but same result. Punting this since it's out of scope.
	 */
	caseInsensitiveRegex, err := regexp.Compile("(?i)" + query)
	if err != nil {
		return nil, fmt.Errorf("Search: %w", err)
	}
	idxs2 := s.SuffixArray.FindAllIndex(caseInsensitiveRegex, -1)
	results := []string{}
	for _, idx := range idxs2 {
		results = append(results, s.CompleteWorks[idx[0]-250:idx[0]+250])
	}
	return results, nil
}
