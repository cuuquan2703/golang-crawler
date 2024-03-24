package service

import (
	"regexp"
	"strings"
)

func LineCount(content string, count chan int) {
	// defer wg.Done()
	count <- strings.Count(content, "\n") + 1
}

func WordCount(content string, count chan int) {
	// defer wg.Done()
	count <- len(strings.Fields((content)))
}

func CharCount(content string, count chan int) {
	// defer wg.Done()
	count <- len(strings.ReplaceAll(content, " ", ""))
}

func Freq(content string, freq chan map[string]int, totalLen chan int) {
	// defer wg.Done()
	c := make(map[string]int)
	length := 0
	reg, _ := regexp.Compile("[^a-zA-Z0-9\\s]+")
	_new := reg.ReplaceAllString(content, "")
	w := strings.Fields(strings.ToLower(_new))
	for _, j := range w {
		length += len(j)
		c[j]++
	}
	freq <- c
	totalLen <- length
}

func Worker(jobs chan string, lineCount chan int, wordCount chan int, charCount chan int, freq chan map[string]int, totalLen chan int) {
	//jobs -> para 10k
	//
	for j := range jobs {
		go LineCount(j, lineCount)
		go WordCount(j, wordCount)
		go CharCount(j, charCount)
		go Freq(j, freq, totalLen)
	}
}
