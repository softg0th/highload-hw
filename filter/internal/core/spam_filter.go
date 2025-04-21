package core

import (
	"filter/internal/entities"
	"math"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

func countStopWords(text string, msgLen int) float64 {
	words := strings.Fields(strings.ToLower(text))

	var stopWordsSum float64

	for _, word := range words {
		word = strings.Trim(word, ".,!?\"'")
		if stopWordWeight, found := entities.StopWords[word]; found {
			stopWordsSum += stopWordWeight
		}
	}

	divisor := math.Log2(float64(msgLen) + 1)

	return stopWordsSum / divisor
}

func parseMsgURL(text string) float64 {
	words := strings.Fields(strings.ToLower(text))

	match, _ := regexp.MatchString(entities.UrlRegexp, words[0])
	if match {
		return 1.0
	}
	return 0.0
}

func checkCapsLock(text string, msgLen int) float64 {
	upperLetters := 0

	for _, letter := range text {
		if unicode.IsUpper(letter) {
			upperLetters++
		}
	}

	return float64(upperLetters) / math.Log2(float64(msgLen)+1)
}

func RunFilterPipeline(msg string) bool {
	msgLen := len(msg)
	var wg sync.WaitGroup
	results := make(chan float64, 3)

	wg.Add(3)

	go func() {
		defer wg.Done()
		results <- countStopWords(msg, msgLen)
	}()

	go func() {
		defer wg.Done()
		results <- parseMsgURL(msg)
	}()

	go func() {
		defer wg.Done()
		results <- checkCapsLock(msg, msgLen)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var sum float64
	var count int

	for v := range results {
		sum += v
		count++
	}

	finalProbability := sum / float64(count)
	return finalProbability >= 0.7
}
