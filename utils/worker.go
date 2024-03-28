package utils

import (
	"regexp"
	"slices"
	"strings"

	"github.com/texttheater/golang-levenshtein/levenshtein"
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

func BoldText(content []string, reqs []string, tag string) []string {
	clone := content
	threshold := 2
	open := tag
	a := strings.Split(tag, " ")[0] + ">"
	close := strings.Replace(a, "<", "</", 1)
	var newContent = make([]string, 0)
	var checked = make([]string, 0)
	for _, el := range clone {
		w := strings.Fields(el)
		for _, word := range w {
			for _, req := range reqs {
				if ((levenshtein.DistanceForStrings([]rune(word), []rune(req), levenshtein.DefaultOptions)) <= threshold) && (!slices.Contains(checked, word)) {
					el = strings.Replace(el, word, open+word+close, -1)
					checked = append(checked, word)
				}
			}
		}
		newContent = append(newContent, el)
	}
	return newContent

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
