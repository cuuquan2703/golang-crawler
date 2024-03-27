package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"webcrawler/logger"
)

var log = logger.CreateLog()

type Data struct {
	data json.RawMessage
}

func Dump(jsonData JSONFile, fileName string) {
	data, err := json.MarshalIndent(jsonData, "", "  ")
	data = bytes.ReplaceAll(data, []byte(`\u003c`), []byte(`<`))
	data = bytes.ReplaceAll(data, []byte(`\u003e`), []byte(`>`))
	if err != nil {
		log.Error("Error in Matshal data ", err)
		return
	}
	err = os.WriteFile(`json/`+fileName, data, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		log.Error("Error writting JSON file ", err)
		return
	}

	log.Info("Write successfully into " + fileName)
}
