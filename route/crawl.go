package route

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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
var limit = 5

func GetAllStatistic(w http.ResponseWriter, r *http.Request) {
	data, err := crawlService.GetAll()
	if err != nil {
		response := &utils.Response{Status: "fail", Message: "cannot get all statistic informations"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &utils.Response{Status: "success", Message: data}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetStatisticByID(w http.ResponseWriter, r *http.Request) {
	_id := r.PathValue("id")
	id, _ := strconv.Atoi(_id)
	data, err := crawlService.GetByID(id)
	if err != nil {
		response := &utils.Response{Status: "fail", Message: "cannot get all statistic informations"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := &utils.Response{Status: "success", Message: data}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CrawlData(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entry")
	body, _ := io.ReadAll(r.Body)
	var b utils.Body
	erJson := json.Unmarshal([]byte(string(body)), &b)
	if erJson != nil {
		response := &utils.Response{Status: "fail", Message: erJson.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	errValid := service.CheckValidURL(b.Url)
	if errValid != nil {
		response := &utils.Response{Status: "fail", Message: errValid.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	matched, _ := service.CheckValidTag(b.Options.BoldTag)
	fmt.Print(b.Options.BoldTag, matched)
	if !matched {
		response := &utils.Response{Status: "fail", Message: "Tag not valid"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if len(b.Url) > limit {
		response := &utils.Response{Status: "fail", Message: "Exceed limit allowed link"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	errCrawl := crawlService.Visit(b.Url, b.Options)

	if errCrawl != nil {
		response := &utils.Response{Status: "fail", Message: errCrawl.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	response := &utils.Response{Status: "success", Message: "Crawling"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Download(w http.ResponseWriter, r *http.Request) {
	fileName := r.PathValue("fileName")
	// Open the JSON file
	file, err := os.Open(fileName)
	if err != nil {
		response := &utils.Response{Status: "fail", Message: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	defer file.Close()
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "application/json")

	_, err = io.Copy(w, file)
	if err != nil {
		response := &utils.Response{Status: "fail", Message: err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
}
