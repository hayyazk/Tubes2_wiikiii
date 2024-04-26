package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Article struct {
	Title string
	URL   string
}

type Result struct {
	Articles          []Article
	ArticleDistance int
	ArticlesVisited   int
	TimeElapsed       string
}

func adjustTitle(str string) string {
	return strings.ReplaceAll(strings.TrimSpace(str), " ", "_")
}

func urlBuilder(title string) string {
	base := "https://en.wikipedia.org/wiki/"
	title = adjustTitle(title)
	return base + title
}

func isAcceptable(url string) (string, bool) {
	if !strings.HasPrefix(url, "/wiki/") || strings.Contains(url, ":") {
		return "https://en.wikipedia.org" + url, false
	}
	return "https://en.wikipedia.org" + url, true
}

func getPath(str string, par map[string]string) []string {
	var path []string
	for par[str] != "none" {
		path = append([]string{str}, path...)
		str = par[str]
	}
	return path
}

func getTitle(link string) string {
	resp, err := http.Get(link)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	title := doc.Find("h1#firstHeading").Text()
	return title
}