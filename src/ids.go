package main

import (
	"fmt"
	"strings"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var depth = make(map[string]int)

func getLinksIDS(current string, parent *map[string]string) []string {
	var urls []string
	d := depth[current] + 1
	resp, err := http.Get(current)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	content := doc.Find("div#bodyContent").Find("a")
	content.Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		// class, _ := item.Attr("class")
		link, valid := isAcceptable(link)
		if !valid {
			return
		}

		// if class == "mw-redirect" {
		// 	title := getTitle(link)
		// 	link = urlBuilder(title)
		// }

		// linkLower := strings.ToLower(link)
		if (*parent)[link] == "" {
			(*parent)[link] = current
			depth[link] = d
			urls = append(urls, link)
		}
	})
	return urls
}

func dfs(stack []string, max int, durl string, parent *map[string]string) ([]string, bool) {
	for len(stack) > 0 {
		if strings.EqualFold(stack[0], durl) {
			return stack, true
		}
		if depth[stack[0]] < max {
			nodes := getLinksIDS(stack[0], parent)
			stack = append(nodes, stack[1:]...)
		} else {
			stack = stack[1:]
		}
	}
	return stack, false
}

func iterSearch(source, dest string, iter int) []string {
	var stack []string
	var parent = make(map[string]string)
	//var vis = make(map[string]bool)

	stack = append(stack, source)
	depth[source] = 0
	parent[source] = "none"
	//vis[strings.ToLower(surl)] = true
	fmt.Println(iter)
	stack, found := dfs(stack, iter, dest, &parent)
	if found {
		return getPath(stack[0], parent)
	} else {
		iter++
		return iterSearch(source, dest, iter)
	}
}

func IDSSearch(source, dest Article) Result {
	var result Result
	//PLACEHOLDER
	start := time.Now()
	result.Articles = append(result.Articles, source)

	path := iterSearch(source.URL, dest.URL, 0)

	for _, p := range path {
		var article Article
		article.URL = p
		article.Title = getTitle(p)
		result.Articles = append(result.Articles, article)
	}

	//result.Articles = append(result.Articles, dest)
	result.ArticlesVisited = 2
	result.ArticleDistance = len(result.Articles) - 1
	result.TimeElapsed = time.Since(start).String()

	return result
}