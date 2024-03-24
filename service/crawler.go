package service

import (
	"fmt"
	"strings"
	"webcrawler/logger"

	"github.com/gocolly/colly"
)

type Crawler struct {
	C *colly.Collector
}

var log = logger.CreateLog()

func (c Crawler) Visit(url string) {
	var para []string
	c.C.OnRequest(func(r *colly.Request) {
		path := strings.Split(r.URL.Path, "-")
		html := path[len(path)-1]
		r.Ctx.Put("fileHTML", html)
		log.Info("Visiting " + r.URL.Path)
	})

	c.C.OnError(func(_ *colly.Response, err error) {
		log.Error("Error: ", err)
	})

	c.C.OnResponse(func(r *colly.Response) {
		filePath := "cache/" + r.Ctx.Get("fileHTML")
		err := r.Save(filePath)
		if err != nil {
			log.Error("Error during saing file ", err)
		} else {
			log.Info("Saved file into " + filePath)
		}
	})

	c.C.OnHTML(".fck_detail", func(e *colly.HTMLElement) {
		e.ForEach("p.Normal:not(:has(script))", func(_ int, kl *colly.HTMLElement) {
			para = append(para, kl.Text)
		})
		fmt.Print(para)
		log.Info("Processing statictis text")
		Concurrency(para)
		log.Info("Done")

	})

	c.C.OnScraped((func(r *colly.Response) {
		log.Info("Finished ")
	}))

	c.C.Visit(url)

}
