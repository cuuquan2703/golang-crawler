package repositories

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"webcrawler/logger"

	_ "github.com/lib/pq"
)

type Para struct {
	Id        int            `json:"id"`
	Url       string         `json:"url"`
	Json      string         `json:"json"`
	LineCount int            `json:"line_count"`
	WordCount int            `json:"word_count"`
	CharCount int            `json:"char_count"`
	AvgLength float64        `json:"avg_length"`
	WordFreq  map[string]int `json:"word_freq"`
}

type ParaRepository struct {
	DB *sql.DB
}

var L = logger.CreateLog()

func ConnectDB() (*sql.DB, error) {
	var url string
	url = "postgresql://postgres:cquan@localhost:5432/web_crawler?sslmode=disable"

	db, err := sql.Open("postgres", url)
	if err != nil {
		L.Error("Error open db:", err)
	}
	return db, err
}

func NewParaRepository() (*ParaRepository, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, err
	}
	return &ParaRepository{
		DB: db,
	}, nil
}

func MapToBase64(dataMap map[string]int) (string, error) {
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		return "", err
	}
	base64Data := base64.StdEncoding.EncodeToString(jsonData)
	return base64Data, err
}

func (repo ParaRepository) Insert(data Para) (sql.Result, error) {
	bytes, err := MapToBase64(data.WordFreq)
	if err != nil {
		fmt.Println(err.Error())
		L.Error("Error in converting map to byets ", err)
	}

	cmd := "INSERT INTO Paragraphs (id,url,json,linecount,wordcount,charcount,wordfreq) values ($1,$2,$3,$4,$5,$6,$7)"
	fmt.Println(len(bytes))
	res, err2 := repo.DB.Exec(cmd, data.Id, data.Url, data.Json, data.LineCount, data.WordCount, data.CharCount, bytes)
	if err2 != nil {
		L.Error("Error: ", err2)
	} else {
		L.Info("Insert successfully")
	}

	return res, err2
}
