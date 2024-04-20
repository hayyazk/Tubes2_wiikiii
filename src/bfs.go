package main

import (
	"time"
)

func BFSSearch(source, dest Article) Result {
	var result Result
	//PLACEHOLDER
	start := time.Now()
	result.Articles = append(result.Articles, source)
	result.Articles = append(result.Articles, dest)
	result.ArticlesVisited = 2
	result.ArticlesTraversed = len(result.Articles)
	result.TimeElapsed = time.Since(start).String()

	return result
}
