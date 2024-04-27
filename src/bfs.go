package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var mut sync.RWMutex

// get links from a page
func getLinksBFS(next chan string, current string, parent *map[string]string) {
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
		link, valid := isAcceptable(link)
		if !valid {
			return
		}
		
		mut.Lock()
		if (*parent)[link] == "" {
			(*parent)[link] = current
			next <- link
		}
		mut.Unlock()
	})
}

// close channel once all waitgroups are done
func closeChan(ch chan string, wg *sync.WaitGroup) {
	(*wg).Wait()
	close(ch)
}

// put list items in channel
func putInChan(arr []string, ch chan string) {
	i := 0
	for _, item := range arr {
		i++
		ch <- item
	}
	close(ch)
}

// scrape multiple pages
func getMultipleLinksBFS(prev, next chan string, wg *sync.WaitGroup, parent *map[string]string) {
	i := 0
	for item := range prev {
		i++
		getLinksBFS(next, item, parent)
	}
	(*wg).Done()
}

func bfs(source string, redirectMap map[string]bool) ([]string, int) {
	// initial declarations
	var parent = make(map[string]string)
	parent[source] = "none"
	var queue []string
	queue = append(queue, source)
	var foundLink string
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
			go getMultipleLinksBFS(prev, next, &wg, &parent)
		}
		go closeChan(next, &wg)

		// put queue items in channel for scraping
		go putInChan(queue, prev)

		// iterate over links recieved from scraping, put in queue
		queue = nil
		fmt.Println(depth)
		for link := range next {
			totalVisited++
			if redirectMap[link] {
				found = true
				foundLink = link
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
		return getPath(foundLink, parent), totalVisited
	}
	return nil, 0
}

func BFSSearch(source, dest Article) Result {
	var result Result
	var redirectMap = make(map[string]bool)
	redirectMap[dest.URL] = true
	getRedirects(buildWhatLinksHereURL(dest.Title), &redirectMap)
	
	start := time.Now()
	result.Articles = append(result.Articles, source)

	path, totalVisited := bfs(source.URL, redirectMap)

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
