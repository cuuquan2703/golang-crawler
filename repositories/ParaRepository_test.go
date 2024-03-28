package repositories_test

import (
	"database/sql"
	"log"
	"reflect"
	"regexp"
	"testing"
	"webcrawler/repositories"
	"webcrawler/utils"

	"github.com/DATA-DOG/go-sqlmock"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}

var db, mock = NewMock()

var repo = repositories.ParaRepository{
	DB: db,
}

func TestInsert(t *testing.T) {
	a, _ := utils.MapToBase64(map[string]int{
		"a": 2,
		"b": 3,
		"c": 4,
	})
	b, _ := utils.Base64ToMap(a)
	newRecord := repositories.ParaModel{
		Id:        123,
		Url:       "/ong-lukashenko-nhung-ke-xa-sung-o-nga-dinh-tron-sang-belarus-4727012.html",
		Json:      "4727012.json",
		LineCount: 12,
		WordCount: 354,
		CharCount: 1234,
		AvgLength: 3.5,
		WordFreq:  a,
	}

	data := repositories.Para{
		Id:        123,
		Url:       "/ong-lukashenko-nhung-ke-xa-sung-o-nga-dinh-tron-sang-belarus-4727012.html",
		Json:      "4727012.json",
		LineCount: 12,
		WordCount: 354,
		CharCount: 1234,
		AvgLength: 3.5,
		WordFreq:  b,
	}
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO Paragraphs (id,url,json,linecount,wordcount,charcount,avglength,wordfreq) values ($1,$2,$3,$4,$5,$6,$7,$8)")).
		WithArgs(newRecord.Id, newRecord.Url, newRecord.Json, newRecord.LineCount, newRecord.WordCount, newRecord.CharCount, newRecord.AvgLength, newRecord.WordFreq).
		WillReturnResult((sqlmock.NewResult(0, 1)))

	_, err := repo.Insert(data)
	if err != nil {
		t.Errorf("Error when Insert data")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetAll(t *testing.T) {
	a, _ := utils.MapToBase64(map[string]int{
		"a": 2,
		"b": 3,
		"c": 4,
	})

	expected := []repositories.ParaModel{{
		Id:        123,
		Url:       "/ong-lukashenko-nhung-ke-xa-sung-o-nga-dinh-tron-sang-belarus-4727012.html",
		Json:      "4727012.json",
		LineCount: 12,
		WordCount: 354,
		CharCount: 1234,
		AvgLength: 3.5,
		WordFreq:  a,
	}}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * from paragraphs")).
		WithoutArgs().
		WillReturnRows(sqlmock.NewRows([]string{"id", "url", "json", "linecount", "wordcount", "charcount", "avglength", "wordfreq"}).
			AddRow(123, "/ong-lukashenko-nhung-ke-xa-sung-o-nga-dinh-tron-sang-belarus-4727012.html", "4727012.json", 12, 354, 1234, 3.5, a))

	res, err := repo.GetAll()
	if err != nil {
		log.Fatal(err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Returned results don't match expected. Expected: %v, Actual: %v", expected, res)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestGetByID(t *testing.T) {
	a, _ := utils.MapToBase64(map[string]int{
		"a": 2,
		"b": 3,
		"c": 4,
	})

	expected := repositories.ParaModel{
		Id:        123,
		Url:       "/ong-lukashenko-nhung-ke-xa-sung-o-nga-dinh-tron-sang-belarus-4727012.html",
		Json:      "4727012.json",
		LineCount: 12,
		WordCount: 354,
		CharCount: 1234,
		AvgLength: 3.5,
		WordFreq:  a,
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * from paragraphs WHERE id=$1")).
		WithArgs(expected.Id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "url", "json", "linecount", "wordcount", "charcount", "avglength", "wordfreq"}).
			AddRow(123, "/ong-lukashenko-nhung-ke-xa-sung-o-nga-dinh-tron-sang-belarus-4727012.html", "4727012.json", 12, 354, 1234, 3.5, a))

	res, err := repo.GetByID(expected.Id)
	if err != nil {
		log.Fatal(err)
	}
	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Returned results don't match expected. Expected: %v, Actual: %v", expected, res)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
