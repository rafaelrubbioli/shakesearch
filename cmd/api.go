package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"pulley.com/shakesearch/pkg/handle"
	"pulley.com/shakesearch/pkg/searcher"
)

func main() {
	searcher, err := searcher.New("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", handle.Search(searcher))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("Listening on port %s...\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
