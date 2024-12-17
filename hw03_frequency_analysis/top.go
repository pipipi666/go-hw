package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(str string) []string {
	wordCount := map[string]int{}
	words := []string{}

	for _, word := range strings.Fields(str) {
		_, ok := wordCount[word]

		if ok {
			wordCount[word]++
		} else {
			wordCount[word] = 1
			words = append(words, word)
		}
	}

	sort.Slice(words, func(i, j int) bool {
		if wordCount[words[i]] == wordCount[words[j]] {
			return words[i] < words[j]
		}

		return wordCount[words[i]] > wordCount[words[j]]
	})

	l := len(words)

	if l >= 10 {
		return words[:10]
	}

	return words[:l]
}
