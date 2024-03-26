package main

import (
	"net/http"
	"webcrawler/route"
)

func main() {
	http.HandleFunc("POST /up", route.CrawlData)
	http.HandleFunc("GET /download/{fileName}", route.Download)

	http.ListenAndServe(":8080", nil)
}
