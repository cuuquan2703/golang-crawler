package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"webcrawler/logger"
	"webcrawler/utils"

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

type ParaModel struct {
	Id        int     `json:"id"`
	Url       string  `json:"url"`
	Json      string  `json:"json"`
	LineCount int     `json:"line_count"`
	WordCount int     `json:"word_count"`
	CharCount int     `json:"char_count"`
	AvgLength float64 `json:"avg_length"`
	WordFreq  string  `json:"word_freq"`
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

func (repo ParaRepository) Insert(data Para) (sql.Result, error) {
	bytes, err := utils.MapToBase64(data.WordFreq)
	if err != nil {
		fmt.Println(err.Error())
		L.Error("Error in converting map to byets ", err)
	}

	cmd := "INSERT INTO Paragraphs (id,url,json,linecount,wordcount,charcount,avglength,wordfreq) values ($1,$2,$3,$4,$5,$6,$7,$8)"
	fmt.Println(len(bytes))
	res, err2 := repo.DB.Exec(cmd, data.Id, data.Url, data.Json, data.LineCount, data.WordCount, data.CharCount, data.AvgLength, bytes)
	if err2 != nil {
		L.Error("Error: ", err2)
	} else {
		L.Info("Insert successfully")
	}

	return res, err2
}

func (repo ParaRepository) GetAll() ([]ParaModel, error) {
	parastatic := []ParaModel{}
	cmd := `SELECT * from paragraphs`
	L.Info("Querying " + cmd)
	row, err := repo.DB.Query(cmd)

	if err != nil {
		L.Error("Error", err)
	} else {
		L.Info("Query successfully")
	}

	for row.Next() {

		para := ParaModel{}
		err := row.Scan(&para.Id, &para.Url, &para.Json, &para.LineCount, &para.WordCount, &para.CharCount, &para.AvgLength, &para.WordFreq)

		if err != nil {
			L.Error("Error", err)
			return nil, err
		}

		parastatic = append(parastatic, para)
	}

	if len(parastatic) == 0 {
		L.Error("Error ", errors.New("no para found"))
		return nil, errors.New("no para found")
	}
	defer row.Close()
	return parastatic, err
}

func (repo ParaRepository) GetByID(id int) (ParaModel, error) {
	para := ParaModel{}
	cmd := `SELECT * from paragraphs WHERE id=$1`
	L.Info("Querying " + cmd)
	row := repo.DB.QueryRow(cmd, id)
	L.Info("Query successfully")
	err := row.Scan(&para.Id, &para.Url, &para.Json, &para.LineCount, &para.WordCount, &para.CharCount, &para.AvgLength, &para.WordFreq)

	if err != nil {
		if err == sql.ErrNoRows {
			L.Error("Error ", errors.New("no para found"))
			return para, errors.New("no para found")
		}
		L.Error("Error ", err)
		return para, errors.New("something went wrong")
	}
	return para, err
}
