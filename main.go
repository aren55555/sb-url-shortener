package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var (
	next            = int64(0)
	idToURL         = map[int64]string{}
	urlToDesination = map[string]string{}
)

func assignID(string) int64 {
	n := next
	next++
	return n
}

func serialize(i int64) string {
	s := strconv.FormatInt(i, 36)
	return fmt.Sprintf("https://short.co/%s", s)
}

func main() {
	http.HandleFunc("/api/generate", func(w http.ResponseWriter, r *http.Request) {
		originalURL := r.URL.Query().Get("url")
		id := assignID(originalURL)
		shortened := serialize(id)

		idToURL[id] = originalURL
		urlToDesination[shortened] = originalURL

		w.Write([]byte(shortened))

		w.Write([]byte("\n\n\n\n"))

		redirAPI, _ := url.Parse("http://localhost:8080/")
		redirAPI.Path = "/api/redirect"
		qp := url.Values{}
		qp.Add("url", shortened)
		redirAPI.RawQuery = qp.Encode()
		w.Write([]byte(redirAPI.String()))
	})

	http.HandleFunc("/api/redirect", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query().Get("url")

		destination, ok := urlToDesination[url]
		if !ok {
			http.Error(w, "not found!", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, destination, http.StatusFound)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
