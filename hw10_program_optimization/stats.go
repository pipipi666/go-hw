package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/valyala/fastjson"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	suff := "." + domain

	for scanner.Scan() {
		line := scanner.Text()
		email := fastjson.GetString([]byte(line), "Email")
		matched := strings.HasSuffix(email, suff)

		if matched {
			num := strings.ToLower(strings.SplitN(email, "@", 2)[1])
			result[num]++
		}
	}

	return result, nil
}
