package service

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"errors"
)

var pattern = `https:\/\/vnexpress\.net\/[^\/]+\.html`
var patternImg = `https:\/\/[^\/]+\.[^\/]+`
var patternTag = `<[^<>]+>`

func CheckValidTag(tag string) (bool, error) {
	var err error
	matched, err := regexp.MatchString(patternTag, tag)
	if matched == false {
		return false, errors.New("Alter tag not valid ")
	}
	return matched, err
}

func CheckValidImgURL(url string) (bool, error) {
	fmt.Println("Check valid url")
	var err error
	matched, err := regexp.MatchString(patternImg, url)
	if matched == false {
		return false, errors.New("URL Img not match")
	}
	return matched, err
}

func CheckValidURL(url []string) error {
	fmt.Println("Check valid url")
	var err error
	for _, i := range url {
		matched, _ := regexp.MatchString(pattern, i)
		if matched == false {
			return errors.New("URL not match")
		}
	}
	return err
}

func CheckCacheURL(url string) bool {
	list := make([]string, 0)
	files, _ := os.ReadDir("./cache")
	for _, f := range files {
		list = append(list, strings.Split(f.Name(), ".")[0])
	}
	if slices.Contains(list, url) {
		return true
	}
	return false
}
