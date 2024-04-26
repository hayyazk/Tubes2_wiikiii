package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// var parent = make(map[string]string)
// var redirect = make(map[string]string)
var mut sync.RWMutex

func getLinksBFS(next chan string, current string, parent *map[string]string) {
	//var links []string
	resp, err := http.Get(current)
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
	content := doc.Find("div#bodyContent").Find("a")
	content.Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		// class, _ := item.Attr("class")
		link, valid := isAcceptable(link)
		if !valid {
			return
		}

		// if class == "mw-redirect" {
		// 	mut.Lock()
		// 	if redirect[link] == "" {
		// 		fmt.Println("hehe")
		// 		title := getTitle(link)
		// 		newlink := urlBuilder(title)
		// 		redirect[link] = newlink
		// 		link = newlink
		// 	} else {
		// 		link = redirect[link]
		// 	}
		// 	mut.Unlock()
		// }
		
		mut.Lock()
		if (*parent)[link] == "" {
			(*parent)[link] = current
			//link = append(urls, link)
			next <- link
		}
		mut.Unlock()
	})
}

func closeChan(ch chan string, wg *sync.WaitGroup) {
	(*wg).Wait()
	close(ch)
}

func putInChan(arr []string, ch chan string) {
	i := 0
	for _, item := range arr {
		i++
		ch <- item
	}
	fmt.Println("Put this amount in prev:", i)
	close(ch)
}

func getMultipleLinks(prev, next chan string, wg *sync.WaitGroup, parent *map[string]string) {
	i := 0
	for item := range prev {
		i++
		getLinksBFS(next, item, parent)
	}
	(*wg).Done()
}

func bfs(source, dest string) ([]string, int) {
	// initial declarations
	var parent = make(map[string]string)
	parent[source] = "none"
	var queue []string
	queue = append(queue, source)
	found := false
	depth := 0
	totalVisited := 1

	// search loop, ends when destination page is found
	for {
		fmt.Println(depth)
		// variables for concurrent search
		var prev = make(chan string)
		var next = make(chan string, 1000)
		var wg sync.WaitGroup

		// setup 100 threads to scrape links
		wg.Add(100)
		for i:=0; i<100; i++ {
			go getMultipleLinks(prev, next, &wg, &parent)
		}
		go closeChan(next, &wg)

		// put queue items in channel for scraping
		go putInChan(queue, prev)

		// iterate over links recieved from scraping, put in queue
		queue = nil
		fmt.Println(depth)
		for link := range next {
			totalVisited++
			if strings.EqualFold(link, dest) {
				found = true
				break
			}
			queue = append(queue, link)
		}
		if found {
			break
		}
		depth++
	}
	if found {
		return getPath(dest, parent), totalVisited
	}
	return nil, 0
}

func BFSSearch(source, dest Article) Result {
	var result Result
	
	start := time.Now()
	result.Articles = append(result.Articles, source)

	path, totalVisited := bfs(source.URL, dest.URL)

	for _, p := range path {
		var article Article
		article.URL = p
		article.Title = getTitle(p)
		result.Articles = append(result.Articles, article)
	}

	result.ArticlesVisited = totalVisited
	result.ArticleDistance = len(result.Articles) - 1
	result.TimeElapsed = time.Since(start).String()

	return result
}
