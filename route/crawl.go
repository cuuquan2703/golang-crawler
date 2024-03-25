package route

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"webcrawler/repositories"
	"webcrawler/service"
	"webcrawler/utils"

	"github.com/gocolly/colly"
)

var a = colly.NewCollector()
var Repo, _ = repositories.NewParaRepository()
var crawlService = service.Crawler{
	C:    a,
	Repo: Repo,
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
