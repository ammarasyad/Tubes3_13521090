package main

import (
	"context"
	"database/sql"
	vector "github.com/niemeyer/golang/src/pkg/container/vector"
	"math"
	"regexp"
	sql_connection "src/backend/src/backend/sql"
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

func ProcessQuestion(question string) string {
	addQuestionRegex := regexp.MustCompile("[Tt]ambahkan pertanyaan [a-z,A-Z,0-9]* dengan jawaban [a-z,A-Z,0-9]*")
	deleteQuestionRegex := regexp.MustCompile("[Hh]apus pertanyaan [a-z,A-Z,0-9]*")
	calendarRegex := regexp.MustCompile("[0-9]{2}/[0-9]{2}/[0-9]{4}")
	calculatorRegex := regexp.MustCompile("[Hh]itung ")
	answer := ""
	if addQuestionRegex.MatchString(question) {

	} else if deleteQuestionRegex.MatchString(question) {

	} else if calendarRegex.MatchString(question) {

	} else if calculatorRegex.MatchString(question) {

	}
	return answer
}

func addQuestion(question string, answer string) {
	datasource := "data source here"
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		panic(err.Error())
	}
	ctx := context.Background()
	conn, err := db.Conn(ctx)
	if err != nil {
		panic(err.Error())
	}
	sql_connection.Create_Question(conn, context.Background(), question, answer)
}

func calculator(expression string) {

}

func infixToPostFix(expression string) string {
	operator := vector.StringVector{}
	ret := ""
	i := 0
	for i < len(expression) {
		s := string(expression[i])
	L1:
		if s == "*" || s == "-" || s == "*" || s == "/" || s == "^" || s == "(" || s == ")" {
			if operator.Len() == 0 {
				operator.Push(string(s))
			} else {
				if s == "(" {
					operator.Push(s)
				} else if s == ")" {
					for operator.Len() != 0 && operator.At(operator.Len()-1) != "(" {
						ret += string(operator.At(operator.Len() - 1))
						operator.Pop()
					}
					if operator.Len() != 0 && operator.At(operator.Len()-1) == "(" {
						operator.Pop()
					}
				} else if getPrecedence(operator.At(operator.Len()-1)) < getPrecedence(s) {
					operator.Push(s)
				} else if getPrecedence(operator.At(operator.Len()-1)) == getPrecedence(s) {
					ret += operator.At(operator.Len() - 1)
					operator.Pop()
					operator.Push(s)
				} else if getPrecedence(operator.At(operator.Len()-1)) > getPrecedence(s) {
					ret += operator.At(operator.Len() - 1)
					operator.Pop()
					goto L1
				}
			}
		} else {
			if s != "(" && s != ")" {
				ret += s
			}
		}
	}
	return ret
}

func getPrecedence(ops string) int {
	switch ops {
	case "+":
	case "-":
		return 1
		break
	case "*":
	case "/":
		return 2
		break
	case "^":
		return 3
		break
	default:
		return 0
	}
	return 0
}
