package main

import (
	"strings"

	"golang.org/x/tour/wc"
)

func WordCount(s string) map[string]int {
	m := make(map[string]int)
	slist := strings.Split(s, " ")
	for _, ch := range slist {
		if val, found := m[ch]; found {
			m[ch] = val + 1
		} else {
			m[ch] = 1
		}
	}
	return m
}

func main() {
	wc.Test(WordCount)
}
