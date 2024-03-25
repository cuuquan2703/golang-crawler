package route

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"webcrawler/service"
	"webcrawler/utils"

	"github.com/gocolly/colly"
)

var a = colly.NewCollector()
var crawlService = service.Crawler{
	C: a,
}

func CrawlData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entry")
	body, _ := io.ReadAll(r.Body)
	var b utils.Body
	er := json.Unmarshal([]byte(string(body)), &b)
	if er != nil {
		fmt.Print("a")
	}
	crawlService.Visit(b.Url, b.Options)
}
