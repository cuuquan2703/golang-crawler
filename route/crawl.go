package route

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"webcrawler/service"

	"github.com/gocolly/colly"
)

var a = colly.NewCollector()
var crawlService = service.Crawler{
	C: a,
}

type Body struct {
	Url string `json:"url"`
}

func CrawlData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entry")
	body, _ := io.ReadAll(r.Body)
	var b Body
	er := json.Unmarshal([]byte(string(body)), &b)
	if er != nil {
		fmt.Print("a")
	}
	fmt.Println(b.Url)
	crawlService.Visit(b.Url)
}
