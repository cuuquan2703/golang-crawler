package service

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"webcrawler/logger"
	"webcrawler/repositories"
	"webcrawler/utils"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

type Crawler struct {
	C    *colly.Collector
	Repo *repositories.ParaRepository
}

var log = logger.CreateLog()

func (c Crawler) Visit(urls []string, options utils.Option) error {

	var err error
	var para []string
	var img []string
	var relatedUrl []string
	var title string
	var parralellism int
	var pattern = `https:\/\/vnexpress\.net\/[^\/]+\.html`
	var re = regexp.MustCompile(pattern)
	c.C.MaxDepth = options.MaxDepth
	c.C.IgnoreRobotsTxt = true
	c.C.AllowURLRevisit = true
	c.C.Async = true
	if len(urls) > 5 {
		parralellism = 5
	} else {
		parralellism = len(urls)
	}
	c.C.Limit(&colly.LimitRule{
		Parallelism: parralellism,
	})
	q, _ := queue.New(len(urls), &queue.InMemoryQueueStorage{MaxSize: 10000})
	c.C.OnRequest(func(r *colly.Request) {

		para = make([]string, 0)
		img = make([]string, 0)
		relatedUrl = make([]string, 0)
		title = ""
		path := strings.Split(r.URL.Path, "-")
		html := path[len(path)-1]
		r.Ctx.Put("url", r.URL.Path)
		r.Ctx.Put("fileHTML", html)
		log.Info("Visiting " + r.URL.Path)
	})

	c.C.OnError(func(e *colly.Response, _err error) {
		err = errors.New(e.Ctx.Get("url") + " " + _err.Error())
		log.Error("Error: ", _err)
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
		fmt.Println(e.Name)
		if slices.Contains(options.Tag, e.Name) || options.Tag[0] == "*" {
			open := `<` + e.Name + `>`
			close := `</` + e.Name + `>`
			title = open + e.Text + close
		}
	})

	c.C.OnHTML("p.description", func(e *colly.HTMLElement) {
		fmt.Println("Element: ", e)
		if slices.Contains(options.Tag, e.Name) || options.Tag[0] == "*" {
			open := `<` + e.Name + `>`
			close := `</` + e.Name + `>`
			para = append(para, open+e.Text+close)
		}
	})

	c.C.OnHTML(".fck_detail", func(e *colly.HTMLElement) {
		e.ForEach("p.Normal:not(:has(script))", func(_ int, kl *colly.HTMLElement) {
			if slices.Contains(options.Tag, kl.Name) || options.Tag[0] == "*" {
				open := `<` + kl.Name + `>`
				close := `</` + kl.Name + `>`
				para = append(para, open+kl.Text+close)
			}
		})
		e.ForEach("a[href]", func(_ int, kl *colly.HTMLElement) {
			if re.MatchString(kl.Attr("href")) {
				relatedUrl = append(relatedUrl, kl.Attr("href"))
			}
		})
		e.ForEach("img[src]", func(_ int, kl *colly.HTMLElement) {
			img = append(img, kl.Attr("src"))
		})
		log.Info("Processing statictis text")
		paras, lineCount, wourdCount, charCount, freq, avgCount := utils.Concurrency(para, options.BoldText)
		//paras, _, _, _, _, _ := Concurrency(para, options.BoldText)

		var json = utils.JSONFile{
			Title:      title,
			Paragraphs: paras,
			ImgUrl:     img,
			RelatedUrl: relatedUrl,
		}

		utils.Dump(json, strings.Split(e.Response.Ctx.Get("fileHTML"), ".")[0]+".json")
		id, _ := strconv.Atoi(strings.Split(e.Response.Ctx.Get("fileHTML"), ".")[0])
		var data = repositories.Para{
			Id:        id,
			Url:       e.Response.Ctx.Get("url"),
			Json:      e.Response.Ctx.Get("fileHTML") + ".json",
			LineCount: lineCount,
			WordCount: wourdCount,
			CharCount: charCount,
			AvgLength: avgCount,
			WordFreq:  freq,
		}

		c.Repo.Insert(data)
		log.Info("Done")

	})

	c.C.OnHTML(".sidebar-1 a[href]", func(e *colly.HTMLElement) {
		if re.MatchString(e.Attr("href")) && !CheckCacheURL(e.Attr("href")) {
			e.Request.Visit(e.Attr("href"))
		}
	})

	c.C.OnScraped((func(r *colly.Response) {
		para = make([]string, 0)
		img = make([]string, 0)
		relatedUrl = make([]string, 0)
		title = ""
		log.Info("Finished ")
	}))

	for _, url := range urls {
		if err := q.AddURL(url); err != nil {
			log.Error("Error adding queue ", err)
		}
	}
	er := q.Run(c.C)
	c.C.Wait()
	if er != nil {
		err = er
		log.Error("Failed to run: ", err)
	}
	return err
}
