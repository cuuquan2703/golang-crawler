package service

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"webcrawler/logger"
	"webcrawler/repositories"
	"webcrawler/utils"

	"github.com/gocolly/colly"
)

type Crawler struct {
	C    *colly.Collector
	Repo *repositories.ParaRepository
}

var log = logger.CreateLog()

func (c Crawler) Visit(url string, options utils.Option, err chan error) {
	var er error
	var para []string
	var img []string
	var relatedUrl []string
	var title string
	var pattern = `https:\/\/vnexpress\.net\/[^\/]+\.html`
	var re = regexp.MustCompile(pattern)
	c.C.MaxDepth = options.MaxDepth
	c.C.Async = true
	c.C.IgnoreRobotsTxt = true
	c.C.OnRequest(func(r *colly.Request) {
		para = make([]string, 0)
		img = make([]string, 0)
		relatedUrl = make([]string, 0)
		path := strings.Split(r.URL.Path, "-")
		html := path[len(path)-1]
		r.Ctx.Put("url", r.URL.Path)
		r.Ctx.Put("fileHTML", html)
		log.Info("Visiting " + r.URL.Path)
	})

	c.C.OnError(func(_ *colly.Response, _err error) {
		er = _err
		log.Error("Error: ", _err)
	})

	c.C.OnResponse(func(r *colly.Response) {
		filePath := "cache/" + r.Ctx.Get("fileHTML")
		err := r.Save(filePath)
		if err != nil {
			er = err
			log.Error("Error during saing file ", err)
		} else {
			log.Info("Saved file into " + filePath)
		}
	})

	c.C.OnHTML(".title-detail", func(e *colly.HTMLElement) {
		fmt.Println(e.Name)
		if e.Name == options.Tag {
			open := `<` + e.Name + `>`
			close := `</` + e.Name + `>`
			title = open + e.Text + close
		}
	})

	c.C.OnHTML("p.description", func(e *colly.HTMLElement) {
		para = append(para, e.Text)
	})

	c.C.OnHTML(".fck_detail", func(e *colly.HTMLElement) {
		e.ForEach("p.Normal:not(:has(script))", func(_ int, kl *colly.HTMLElement) {
			para = append(para, kl.Text)
		})
		e.ForEach("a[href]", func(_ int, kl *colly.HTMLElement) {
			if re.MatchString(kl.Attr("href")) {
				relatedUrl = append(relatedUrl, kl.Attr("href"))
			}
		})
		e.ForEach("img[src]", func(_ int, kl *colly.HTMLElement) {
			img = append(img, kl.Attr("src"))
		})
		fmt.Print(para)
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
		id, err := strconv.Atoi(strings.Split(e.Response.Ctx.Get("fileHTML"), ".")[0])
		er = err
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
		if re.MatchString(e.Attr("href")) {
			e.Request.Visit(e.Attr("href"))
		}
	})

	c.C.OnScraped((func(r *colly.Response) {
		log.Info("Finished ")
	}))

	c.C.Visit(url)
	c.C.Wait()
	err <- er
}
