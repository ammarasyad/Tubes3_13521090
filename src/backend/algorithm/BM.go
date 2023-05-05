package algorithm

import "math"
import "strings"

func BM(text string, pattern string) bool {
	text = strings.ToLower(text)
	pattern = strings.ToLower(pattern)
	last := buildLast(pattern)
	i := len(pattern) - 1
	if i > len(text)-1 {
		return false
	}
	j := len(pattern) - 1
	for i <= len(text)-1 {
		if pattern[j] == text[i] {
			if j == 0 {
				return true
			} else {
				i--
				j--
			}
		} else {
			lo := last[text[i]]
			i = i + len(pattern) - int(math.Min(float64(j), float64(1+lo)))
			j = len(pattern) - 1
		}
	}
	return false
}

func buildLast(pattern string) []int {
	last := make([]int, 128)
	for i := 0; i < len(last); i++ {
		last[i] = -1
	}
	for i := 0; i < len(pattern); i++ {
		last[int(pattern[i])] = i
	}
	return last
}
