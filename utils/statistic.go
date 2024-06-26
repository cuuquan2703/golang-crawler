package utils

import "fmt"

func Concurrency(list []string, boldText []string, tag string) ([]string, int, int, int, map[string]int, float64) {
	numWorker := 5
	numJobs := len(list)
	jobs := make(chan string)
	lineCount := make(chan int)
	wordCount := make(chan int)
	charCount := make(chan int)
	freq := make(chan map[string]int)
	totalLen := make(chan int)

	lc, wc, cc, fr, ttl := 0, 0, 0, make(map[string]int), 0
	for w := 1; w <= numWorker; w++ {
		go Worker(jobs, lineCount, wordCount, charCount, freq, totalLen)
	}

	for _, para := range list {
		jobs <- para
	}
	close(jobs)

	for res := 1; res <= numJobs; res++ {
		_lc := <-lineCount
		_wc := <-wordCount
		_cc := <-charCount
		_fr := <-freq
		_ttl := <-totalLen
		lc += _lc
		wc += _wc
		cc += _cc
		ttl += _ttl
		for key, value := range _fr {
			fr[key] += value
		}

	}
	avgCount := float64(ttl) / float64(wc)
	fmt.Print(lc, wc, cc, avgCount, freq)
	nc := BoldText(list, boldText, tag)
	return nc, lc, wc, cc, fr, avgCount
}
