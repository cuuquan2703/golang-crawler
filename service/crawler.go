package service

import (
	"fmt"
	"strings"
	"webcrawler/logger"
	"webcrawler/utils"

	"github.com/gocolly/colly"
)

type Crawler struct {
	C *colly.Collector
}

var log = logger.CreateLog()

func (c Crawler) Visit(url string, options utils.Option) {
	var para []string
	var img []string
	var relatedUrl []string
	var title string
	c.C.MaxDepth = options.MaxDepth
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

	c.C.OnHTML(".title-detail", func(e *colly.HTMLElement) {
		title = e.Text
	})

	c.C.OnHTML("p.description", func(e *colly.HTMLElement) {
		para = append(para, e.Text)
	})

	c.C.OnHTML(".fck_detail", func(e *colly.HTMLElement) {
		e.ForEach("p.Normal:not(:has(script))", func(_ int, kl *colly.HTMLElement) {
			para = append(para, kl.Text)
		})
		e.ForEach("a[href]", func(_ int, kl *colly.HTMLElement) {
			relatedUrl = append(relatedUrl, kl.Attr("href"))
		})
		e.ForEach("img[src]", func(_ int, kl *colly.HTMLElement) {
			img = append(img, kl.Attr("src"))
		})
		fmt.Print(para)
		log.Info("Processing statictis text")
		// paras, lineCount, wourdCount, charCount, freq, avgCount := Concurrency(para, options.BoldText)
		paras, _, _, _, _, _ := Concurrency(para, options.BoldText)

		var json = utils.JSONFile{
			Title:      title,
			Paragraphs: paras,
			ImgUrl:     img,
			RelatedUrl: relatedUrl,
		}

		Dump(json, e.Response.Ctx.Get("fileHTML")+".json")

		log.Info("Done")

	})

	c.C.OnScraped((func(r *colly.Response) {
		log.Info("Finished ")
	}))

	c.C.Visit(url)

}
