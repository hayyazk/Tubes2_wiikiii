package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"time"
)

var t = template.Must(template.ParseFiles(path.Join("view", "index.html")))

func mainHandler(w http.ResponseWriter, r *http.Request) {
	var result Result
	result.ArticleDistance = -2
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
	tempURL := urlBuilder(sourceTitle)
	source.Title = getTitle(tempURL)
	source.URL = urlBuilder(source.Title)

	var dest Article
	tempURL = urlBuilder(destTitle)
	dest.Title = getTitle(tempURL)
	dest.URL = urlBuilder(dest.Title)

	var result Result

	if dest.Title == source.Title {
		start := time.Now()
		result.ArticleDistance = 0
		result.ArticlesVisited = 1
		result.Articles = append(result.Articles, dest)
		result.TimeElapsed = time.Since(start).Milliseconds()
	} else {
		if algorithm == "bfs" {
			result = BFSSearch(source, dest)
		}
		if algorithm == "ids" {
			result = IDSSearch(source, dest)
		}
	}

	t.Execute(w, result)
}

func main() {
	mux := http.NewServeMux()
	fmt.Println("http://localhost:8080/")
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("view/stylesheet"))))
	mux.HandleFunc("/", mainHandler)
	mux.HandleFunc("/search", searchHandler)
	http.ListenAndServe(":8080", mux)
}
