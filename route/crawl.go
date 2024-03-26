package route

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	err := make(chan error)
	fmt.Println("entry")
	body, _ := io.ReadAll(r.Body)
	var b utils.Body
	er := json.Unmarshal([]byte(string(body)), &b)
	if er != nil {
		fmt.Print("a")
	}
	go crawlService.Visit(b.Url, b.Options, err)
	_err := <-err
	if _err != nil {
		response := &utils.Response{Status: "fail", Message: _err.Error()}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	response := &utils.Response{Status: "success", Message: ""}
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
