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
	"strconv"
)

const PreviewSize = 250
const PageSize = 20

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
		query, qOK := r.URL.Query()["q"]
		if !qOK || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}

		offset, offsetExists := r.URL.Query()["offset"]
		offsetAsInt := 0
		if offsetExists {
			parsedInt, parseIntError := strconv.Atoi(offset[0])
			if parseIntError != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("invalid offset in URL params"))
				return
			}
			offsetAsInt = parsedInt
		}

		results, hasMore, searchError := searcher.Search(query[0], offsetAsInt)
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
		w.Header().Set("X-Has-More", strconv.FormatBool(hasMore))
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
func (s *Searcher) Search(query string, offset int) ([]string, bool, error) {
	/**
	 * TODO: Something's wrong with some punctuations. Hunch is the \r\n characters that got in the
	 * way.
	 *
	 * - period doesn't work -- "posting is no need. " works; "posting is no need. O" doesn't.
	 * - comma works-- "very finely, very comely".
	 */
	caseInsensitiveRegex, err := regexp.Compile("(?i)" + query)
	if err != nil {
		return nil, false, fmt.Errorf("Search: %w", err)
	}

	idxs := s.SuffixArray.FindAllIndex(caseInsensitiveRegex, -1)

	totalCount := len(idxs)
	boundedPageStartIndex := Min(offset, totalCount)
	boundedPageEndIndex := Min(boundedPageStartIndex+PageSize, totalCount)

	results := []string{}
	for _, idx := range idxs[boundedPageStartIndex:boundedPageEndIndex] {
		boundedPreviewStartIndex := Max(idx[0]-PreviewSize, 0)
		boundedPreviewEndIndex := Min(idx[0]+PreviewSize, len(s.CompleteWorks))
		results = append(results, s.CompleteWorks[boundedPreviewStartIndex:boundedPreviewEndIndex])
	}

	return results, boundedPageEndIndex < totalCount, nil
}

// Min returns the smaller of x or y.
func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// Max returns the larger of x or y.
func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
