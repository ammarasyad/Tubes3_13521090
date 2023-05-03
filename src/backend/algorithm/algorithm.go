package main

import (
	"math"
)

func KMP(text string, pattern string) bool {
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

func BM(text string, pattern string) bool {
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

func levenshteinDistance(text string, pattern string) int {
	d := make([][]int, len(text)+1)
	for i := 0; i < len(text)+1; i++ {
		d[i] = make([]int, len(pattern)+1)
		for j := 0; j < len(pattern)+1; j++ {
			d[i][j] = 0
			if i == 0 {
				d[i][j] = j
			} else if j == 0 {
				d[i][j] = i
			}
		}
	}
	cost := 0
	for i := 1; i < len(d); i++ {
		for j := 1; j < len(d[i]); j++ {
			if pattern[j-1] == text[i-1] {
				cost = 0
			} else {
				cost = 1
			}
			d[i][j] = int(math.Min(math.Min(float64(d[i-1][j]+1), float64(d[i][j-1]+1)), float64(d[i-1][j-1]+cost)))
		}
	}
	return d[len(text)][len(pattern)]
}
