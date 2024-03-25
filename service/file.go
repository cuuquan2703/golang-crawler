package service

import (
	"encoding/json"
	"fmt"
	"os"
	"webcrawler/utils"
)

func Dump(jsonData utils.JSONFile, fileName string) {
	data, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		log.Error("Error in Matshal data ", err)
		return
	}
	err = os.WriteFile(fileName, data, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		log.Error("Error writting JSON file ", err)
		return
	}

	log.Info("Write successfully into " + fileName)
}
