package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

type Article struct {
	Title string
	URL   string
}

type Result struct {
	Articles          []Article
	ArticlesTraversed int
	ArticlesVisited   int
	TimeElapsed       string
}

var t = template.Must(template.ParseFiles("index.html"))

func urlBuilder(title string) string {
	base := "https://en.wikipedia.org/wiki/"
	cleanTitle := strings.ReplaceAll(strings.TrimSpace(title), " ", "_")
	return base + cleanTitle
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	var result Result
	result.ArticlesTraversed = 0
	t.Execute(w, result)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		fmt.Println(err)
		return
	}

	queries := u.Query()
	sourceTitle := queries.Get("source")
	destTitle := queries.Get("dest")
	algorithm := queries.Get("algorithm")

	var source Article
	source.Title = sourceTitle
	source.URL = urlBuilder(sourceTitle)

	var dest Article
	dest.Title = destTitle
	dest.URL = urlBuilder(destTitle)

	var result Result
	if algorithm == "bfs" {
		result = BFSSearch(source, dest)
	}
	if algorithm == "ids" {
		result = IDSSearch(source, dest)
	}

	t.Execute(w, result)
}

func main() {
	mux := http.NewServeMux()
	fmt.Println("http://localhost:8080/")
	mux.HandleFunc("/", mainHandler)
	mux.HandleFunc("/search", searchHandler)
	http.ListenAndServe(":8080", mux)
}
