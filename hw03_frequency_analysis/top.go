package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var reg = regexp.MustCompile("^[[:punct:]]+|[[:punct:]]+$")

func Top10(str string) []string {
	wordCount := map[string]int{}
	words := []string{}

	for _, word := range strings.Fields(str) {
		formattedWord := reg.ReplaceAllString(strings.ToLower(word), "")

		if len(formattedWord) == 0 {
			continue
		}

		_, ok := wordCount[formattedWord]

		if ok {
			wordCount[formattedWord]++
		} else {
			wordCount[formattedWord] = 1
			words = append(words, formattedWord)
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
