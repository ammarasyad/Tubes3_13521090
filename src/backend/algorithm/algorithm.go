package algorithm

import (
	"database/sql"
	"errors"
	"math"
	"regexp"
	"sort"
	sql_connection "src/backend/src/backend/sql"
	"strconv"
	"strings"
	"time"

	vector "github.com/niemeyer/golang/src/pkg/container/vector"
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

func ProcessQuestion(db *sql.DB, question string, kmpbm bool) string {
	addQuestionRegex := regexp.MustCompile("[Tt]ambahkan pertanyaan \\b[^.?!]+ dengan jawaban \\b[^.?!]+")
	deleteQuestionRegex := regexp.MustCompile("[Hh]apus pertanyaan [a-z,A-Z,0-9]*")
	calendarRegex := regexp.MustCompile("[0-9]{2}/[0-9]{2}/[0-9]{4}")
	calculatorRegex := regexp.MustCompile("[Hh]itung ")
	answer := ""
	if addQuestionRegex.MatchString(question) {
		temp := strings.Replace(question, "Tambahkan pertanyaan ", "", 1)
		temp = strings.Replace(temp, "tambahkan pertanyaan ", "", 1)
		temp = strings.Replace(temp, " dengan jawaban ", "|", 1)
		strarr := strings.Split(temp, "|")
		answer = addQuestion(db, strarr[0], strarr[1], kmpbm)
	} else if deleteQuestionRegex.MatchString(question) {
		temp := strings.Replace(question, "Hapus pertanyaan ", "", 1)
		temp = strings.Replace(temp, "hapus pertanyaan ", "", 1)
		answer = deleteQuestion(db, temp, kmpbm)
	} else if calendarRegex.MatchString(question) {
		//temp := strings.Replace(question, "/", "-", -1)
		// temp := strings.ReplaceAll()
		date, err := time.Parse("02/01/2006", question)
		if err != nil {
			return err.Error()
		}
		return date.Weekday().String()
	} else if calculatorRegex.MatchString(question) {
		temp := strings.Replace(question, "Hitung ", "", 1)
		temp = strings.Replace(temp, "hitung ", "", 1)
		answer = calculator(temp)
	} else {
		answer = answerQuestion(db, question, kmpbm)
	}
	return answer
}

func addQuestion(db *sql.DB, question string, answer string, kmpbm bool) string {
	questions := sql_connection.Read_All_Questions(db)
	message := ""
	for i := 0; i < len(questions); i++ {
		if len(questions[i]) == len(question) {
			if kmpbm {
				if KMP(questions[i], question) {
					message = "partanyaan " + question + " sudah ada! jawaban diupdate ke " + answer
					sql_connection.Update_Answer(db, question, answer)
					return message
				}
			} else {
				if BM(questions[i], question) {
					message = "partanyaan " + question + " sudah ada! jawaban diupdate ke " + answer
					sql_connection.Update_Answer(db, question, answer)
					return message
				}
			}
		}
	}
	sql_connection.Create_Question(db, question, answer)
	message = "pertanyaan " + question + " telah ditambah"
	return message
}

func deleteQuestion(db *sql.DB, question string, kmpbm bool) string {
	questions := sql_connection.Read_All_Questions(db)
	message := ""
	for i := 0; i < len(questions); i++ {
		if len(questions[i]) == len(question) {
			if kmpbm {
				if KMP(questions[i], question) {
					message = "partanyaan " + question + " telah dihapus"
					sql_connection.Delete_Question(db, question)
					return message
				}
			} else {
				if BM(questions[i], question) {
					message = "partanyaan " + question + " telah dihapus"
					sql_connection.Delete_Question(db, question)
					return message
				}
			}
		}
	}
	message = "tidak ada pertanyaan " + question + " di database"
	return message
}

func answerQuestion(db *sql.DB, question string, kmpbm bool) string {
	questions := sql_connection.Read_All_Questions(db)
	message := ""
	if len(questions) == 0 {
		return "tidak ada pertanyaan di database"
	} else {
		for i := 0; i < len(questions); i++ {
			if len(question) == len(questions[i]) {
				if kmpbm {
					if KMP(questions[i], question) {
						message = sql_connection.Read_Question(db, questions[i])
						return message
					}
				} else {
					if BM(questions[i], question) {
						message = sql_connection.Read_Question(db, questions[i])
						return message
					}
				}
			}

		}
		sort.SliceStable(questions, func(i, j int) bool {
			return levenshteinDistance(questions[i], question) < levenshteinDistance(questions[j], question)
		})
		percentage := ((float64(len(question)) - float64(levenshteinDistance(questions[0], question))) / float64(len(question))) * 100
		if percentage >= 90 {
			message = sql_connection.Read_Question(db, questions[0])
			return message
		}
		message = "tidak ada pertanyaan " + question + " di database.\n" +
			"apakah maksud anda: \n"
		i := 0
		for i < len(questions) && i < 3 {
			message += strconv.Itoa(i+1) + ". " + questions[i] + "\n"
			i++
		}
		return message
	}

}

func calculator(expression string) string {
	ret, err := calculatePostfix(infixToPostFix(expression))
	if err != nil {
		return err.Error()
	}
	return strconv.FormatFloat(ret, 'f', 3, 32)
}

func infixToPostFix(expression string) string {
	operator := vector.StringVector{}
	ret := ""
	i := 0
	for i < len(expression) {
		s := string(expression[i])
	L1:
		if s == "+" || s == "-" || s == "*" || s == "/" || s == "^" || s == "(" || s == ")" {
			if operator.Len() == 0 {
				operator.Push(s)
			} else {
				if s == "(" {
					operator.Push(s)
				} else if s == ")" {
					for operator.Len() > 0 && operator.At(operator.Len()-1) != "(" {
						ret += operator.At(operator.Len() - 1)
						operator.Pop()
					}
					if operator.Len() > 0 && operator.At(operator.Len()-1) == "(" {
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
		i++
	}
	for operator.Len() > 0 {
		ret += operator.At(operator.Len() - 1)
		operator.Pop()
	}
	return ret
}

func getPrecedence(ops string) int {
	if ops == "-" || ops == "+" {
		return 1
	} else if ops == "*" || ops == "/" {
		return 2
	} else if ops == "^" {
		return 3
	} else {
		return 0
	}
}

func calculatePostfix(expression string) (float64, error) {
	ret := vector.Vector{}
	temp1 := 0.
	temp2 := 0.
	for i := 0; i < len(expression); i++ {
		s := string(expression[i])
		if s == "+" || s == "-" || s == "*" || s == "/" || s == "^" {
			if ret.Len() < 2 {
				return 0, errors.New("unexpected expression")
			}
			temp2 = ret.At(ret.Len() - 1).(float64)
			ret.Pop()
			temp1 = ret.At(ret.Len() - 1).(float64)
			ret.Pop()
			if s == "+" {
				res := temp1 + temp2
				ret.Push(res)
			} else if s == "-" {
				res := temp1 - temp2
				ret.Push(res)
			} else if s == "*" {
				res := temp1 * temp2
				ret.Push(res)
			} else if s == "/" {
				if temp2 == 0 {
					//panic(errors.New("zero division error"))
					return 0, errors.New("zero division error")
				}
				res := temp1 / temp2
				ret.Push(res)
			} else if s == "^" {
				res := math.Pow(temp1, temp2)
				ret.Push(res)
			} else {
				//panic(errors.New("unexpected operand"))
				return 0, errors.New("unexpected operand")
			}
		} else {
			res, err := strconv.ParseFloat(s, 64)
			if err != nil {
				panic(err)
			}
			ret.Push(res)
		}
	}
	return ret.At(0).(float64), nil
}
