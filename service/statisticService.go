package service

import (
	"webcrawler/repositories"
	"webcrawler/utils"
)

func (c Crawler) GetAll() ([]repositories.Para, error) {
	parasModel, err := c.Repo.GetAll()
	var paras []repositories.Para

	for i := range parasModel {
		para := repositories.Para{}
		para.Id = parasModel[i].Id
		para.Json = parasModel[i].Json
		para.Url = parasModel[i].Url
		para.LineCount = parasModel[i].LineCount
		para.WordCount = parasModel[i].WordCount
		para.CharCount = parasModel[i].CharCount
		para.AvgLength = parasModel[i].AvgLength
		para.WordFreq, _ = utils.Base64ToMap(parasModel[i].WordFreq)
		paras = append(paras, para)
	}
	return paras, err
}

func (c Crawler) GetByID(id int) (repositories.Para, error) {
	parasModel, err := c.Repo.GetByID(id)
	para := repositories.Para{}
	para.Id = parasModel.Id
	para.Json = parasModel.Json
	para.Url = parasModel.Url
	para.LineCount = parasModel.LineCount
	para.WordCount = parasModel.WordCount
	para.CharCount = parasModel.CharCount
	para.AvgLength = parasModel.AvgLength
	para.WordFreq, _ = utils.Base64ToMap(parasModel.WordFreq)

	return para, err
}
