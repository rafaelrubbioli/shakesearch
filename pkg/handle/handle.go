package handle

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"pulley.com/shakesearch/pkg/searcher"
)

func Search(searcher searcher.Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			_, err := w.Write([]byte("missing search query in URL params"))
			if err != nil {
				log.Println(err)
			}

			return
		}

		results := make([]string, 0)
		for _, q := range query {
			results = append(results, searcher.Search(q)...)
		}

		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte("encoding failure"))
			if err != nil {
				log.Println(err)
			}

			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(buf.Bytes())
		if err != nil {
			log.Println(err)
		}
	}
}
