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

func buildWhatLinksHereURL(title string) string {
	return "https://en.wikipedia.org/wiki/Special:WhatLinksHere?target=" + adjustTitle(title) + "&namespace=0&limit=500&hidetrans=1&hidelinks=1"
}

func getRedirects(link string, redirectMap *map[string]bool) {
	resp, err := http.Get(link)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	content := doc.Find("ul#mw-whatlinkshere-list").Find("a.mw-redirect")
	content.Each(func(index int, item *goquery.Selection) {
		title := item.Text();
		if title != "edit" {
			(*redirectMap)[urlBuilder(title)] = true
		}
	})
	next, exist := doc.Find("div.mw-pager-navigation-bar").Find("a.mw-nextlink").Attr("href")
	if exist {
		getRedirects("https://en.wikipedia.org" + next, redirectMap)
	} 
}