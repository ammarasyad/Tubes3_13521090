package algorithm

import "strings"

func KMP(text string, pattern string) bool {
	text = strings.ToLower(text)
	pattern = strings.ToLower(pattern)
	fail := computeBorder(pattern)
	i := 0
	j := 0
	for i < len(text) {
		if pattern[j] == text[i] {
			if j == len(pattern)-1 {
				return true
			}
			i++
			j++
		} else if j > 0 {
			j = fail[j-1]
		} else {
			i++
		}
	}
	return false
}

func computeBorder(pattern string) []int {
	fail := make([]int, len(pattern))
	for i := 0; i < len(fail); i++ {
		fail[i] = 0
	}
	j := 0
	i := 1
	for i < len(pattern) {
		if pattern[j] == pattern[i] {
			fail[i] = j + 1
			i++
			j++
		} else if j > 0 {
			j = fail[j-1]
		} else {
			fail[i] = 0
			i++
		}
	}
	return fail
}
