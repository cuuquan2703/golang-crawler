package main

import (
	"net/http"
	"webcrawler/route"
)

func main() {
	http.HandleFunc("GET /up", route.CrawlData)

	http.ListenAndServe(":8080", nil)
}
