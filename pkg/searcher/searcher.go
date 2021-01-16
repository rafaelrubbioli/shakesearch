//go:generate go run -mod=mod github.com/golang/mock/mockgen -package=mocks -source=$GOFILE -destination=../../mocks/searcher.go

package searcher

import (
	"index/suffixarray"
	"io/ioutil"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
)

var endingCharacters = map[string]struct{}{
	".":    {},
	"?":    {},
	"!":    {},
	"\r\n": {},
}

type Searcher interface {
	Search(query string) []string
}

func New(fileName string) (Searcher, error) {
	searcher := searcher{}
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	searcher.CompleteWorks = string(dat)
	searcher.SuffixArray = suffixarray.New([]byte(strings.ToLower(string(dat))))
	searcher.cache = cache.New(time.Hour, time.Hour)

	return searcher, nil
}

type searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
	cache         *cache.Cache
}

func (s searcher) Search(query string) []string {
	query = strings.ToLower(query)
	cached, ok := s.cache.Get(query)
	if ok {
		return cached.([]string)
	}

	indexes := s.SuffixArray.Lookup([]byte(query), -1)
	results := make([]string, 0)
	for _, idx := range indexes {
		start := s.findBeginning(idx)
		end := s.findEnd(idx)
		results = append(results, `"`+s.CompleteWorks[start:end]+`"`)
	}

	s.cache.Set(query, results, time.Hour)
	return results
}

func (s searcher) findBeginning(index int) int {
	index = index - 1
	for index > 0 {
		if _, ok := endingCharacters[string(s.CompleteWorks[index])]; ok {
			return index + 1
		}

		index = index - 1
	}

	return 0
}

func (s searcher) findEnd(index int) int {
	index = index + 1
	for index < len(s.CompleteWorks) {
		if _, ok := endingCharacters[string(s.CompleteWorks[index])]; ok {
			return index + 1
		}

		index = index + 1
	}

	return len(s.CompleteWorks)
}
