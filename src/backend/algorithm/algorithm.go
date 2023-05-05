package algorithm

import (
	"database/sql"
	"math"
	"regexp"
	"sort"
	sql_connection "src/backend/src/backend/sql"
	"strconv"
	"strings"
	"time"
)

func levenshteinDistance(text string, pattern string) int {
	text = strings.ToLower(text)
	pattern = strings.ToLower(pattern)
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
	addQuestionRegex := regexp.MustCompile(`[Tt]ambahkan pertanyaan \b[^.?!]+ dengan jawaban \b[^.?!]+`)
	deleteQuestionRegex := regexp.MustCompile(`[Hh]apus pertanyaan \b[^.?!]+`)
	calendarRegex := regexp.MustCompile(`^([Hh]ari apa\s)?([0-9]{2}\/[0-9]{2}\/[0-9]{4})$`)
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
		temp := strings.Replace(question, "Hari apa ", "", 1)
		temp = strings.Replace(temp, "hari apa ", "", 1)
		date, err := time.Parse("02/01/2006", temp)
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
					message = "pertanyaan " + question + " sudah ada! jawaban diupdate menjadi " + answer
					sql_connection.Update_Answer(db, questions[i], answer)
					return message
				}
			} else {
				if BM(questions[i], question) {
					message = "pertanyaan " + question + " sudah ada! jawaban diupdate menjadi " + answer
					sql_connection.Update_Answer(db, questions[i], answer)
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
					message = "pertanyaan " + question + " telah dihapus"
					sql_connection.Delete_Question(db, questions[i])
					return message
				}
			} else {
				if BM(questions[i], question) {
					message = "pertanyaan " + question + " telah dihapus"
					sql_connection.Delete_Question(db, questions[i])
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
