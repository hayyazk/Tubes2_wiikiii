package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var depth = make(map[string]int)
// var childOfDepth = make(map[int][]string)
// var mtx sync.RWMutex


// get links
func getLinksIDS(current string, parent *map[string]string) []string {
	var urls []string
	d := depth[current] + 1
	resp, err := http.Get(current)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	content := doc.Find("div#bodyContent").Find("a")
	content.Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		link, valid := isAcceptable(link)
		if !valid {
			return
		}
		if (*parent)[link] == "" {
			(*parent)[link] = current
			depth[link] = d
			urls = append(urls, link)
		}
	})
	return urls
}

// // get links but to a channel
// func getLinksIDSConcurrent(next chan string, current string, parent *map[string]string) {
// 	// fmt.Println("soo...")
// 	resp, err := http.Get(current)
// 	if err != nil {
// 		//fmt.Println(err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	doc, err := goquery.NewDocumentFromReader(resp.Body)
// 	if err != nil {
// 		//fmt.Println(err)
// 		return
// 	}

// 	content := doc.Find("div#bodyContent").Find("a")
// 	content.Each(func(index int, item *goquery.Selection) {
// 		link, _ := item.Attr("href")
// 		link, valid := isAcceptable(link)
// 		if !valid {
// 			return
// 		}
// 		mtx1.Lock()
// 		if (*parent)[link] == "" {
// 			(*parent)[link] = current
// 			depth[link] = depth[current] + 1
// 			next <- link
// 		}
// 		mtx1.Unlock()
// 	})
// }

// func getMultipleLinksIDS(prev, next chan string, wg *sync.WaitGroup, parent *map[string]string) {
// 	i := 0
// 	for item := range prev {
// 		i++
// 		getLinksIDSConcurrent(next, item, parent)
// 	}
// 	(*wg).Done()
// }

// func BFSLike(stack []string, max int, redirectMap map[string]bool, parent *map[string]string) []string {
// 	for len(stack) > 0 && depth[stack[0]] == max {
// 		// fmt.Println("IS IT HERE")
// 		if redirectMap[stack[0]] {
// 			return stack
// 		}
// 		stack = stack[1:]
// 	}

// 	if len(stack) == 0 {
// 		return nil
// 	}
// 	var prev = make(chan string)
// 	var next = make(chan string, 1000)
// 	var wg sync.WaitGroup

// 	wg.Add(80)
// 	for i:=0; i<80; i++ {
// 		go getMultipleLinksIDS(prev, next, &wg, parent)
// 	}
// 	go closeChan(next, &wg)
// 	go putInChan(stack, prev)

// 	stack = nil
// 	for link := range next {
// 		// fmt.Println("heyyy")
// 		stack = append([]string{link}, stack...)
// 	}
// 	fmt.Println("eyy")

// 	return stack
// }

// depth limited search
func dls(stack []string, max int, redirectMap map[string]bool, parent *map[string]string) ([]string, bool, int) {
	totalVisited := 0
	for len(stack) > 0 {
		totalVisited++
		if redirectMap[stack[0]] {
			return stack, true, totalVisited
		}

		if depth[stack[0]] < max {
			nodes := getLinksIDS(stack[0], parent)
			stack = append(nodes, stack[1:]...)
		} else {
			stack = stack[1:]
		}
	}
	return stack, false, totalVisited
}

// handle iterations of ids
func iterate(source string, redirectMap map[string]bool, iter int) ([]string, int) {
	var stack []string
	var parent = make(map[string]string)

	stack = append(stack, source)
	depth[source] = 0
	parent[source] = "none"

	// fmt.Println(iter)
	stack, found, totalVisited := dls(stack, iter, redirectMap, &parent)
	if found {
		return getPath(stack[0], parent), totalVisited
	} else {
		iter++
		return iterate(source, redirectMap, iter)
	}
}

func IDSSearch(source, dest Article) Result {
	var result Result
	var redirectMap = make(map[string]bool)
	redirectMap[dest.URL] = true
	getRedirects(buildWhatLinksHereURL(dest.Title), &redirectMap)
	
	start := time.Now()
	result.Articles = append(result.Articles, source)

	path, totalVisited := iterate(source.URL, redirectMap, 0)

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