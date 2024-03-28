package utils

import (
	"encoding/base64"
	"encoding/json"
)

func MapToBase64(dataMap map[string]int) (string, error) {
	jsonData, err := json.Marshal(dataMap)
	if err != nil {
		return "", err
	}
	base64Data := base64.StdEncoding.EncodeToString(jsonData)
	return base64Data, err
}

func Base64ToMap(data string) (map[string]int, error) {
	var err error
	jsonData, errEnc := base64.StdEncoding.DecodeString(data)
	if errEnc != nil {
		err = errEnc
	}
	var dataMap map[string]int
	json.Unmarshal(jsonData, &dataMap)

	return dataMap, err
}
