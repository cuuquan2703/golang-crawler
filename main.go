package main

import (
	"net/http"
	"webcrawler/route"
)

func main() {
	http.HandleFunc("POST /up", route.CrawlData)
	http.HandleFunc("GET /download/{fileName}", route.Download)
	http.HandleFunc("GET /statistic/all", route.GetAllStatistic)
	http.HandleFunc("GET /statistic/{id}", route.GetStatisticByID)

	http.ListenAndServe(":8080", nil)
}
