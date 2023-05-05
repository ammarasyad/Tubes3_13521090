package sql_connection

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// CRD operations for history table (no update)

func Create_History(db *sql.DB, question string, answer string) {
	_, err := db.Exec("INSERT IGNORE INTO history (question, answer) VALUES (?, ?)", question, answer)

	if err != nil {
		panic(err.Error())
	}
}

func Read_All_History(db *sql.DB) ([]string, []string) {
	var questions []string
	var answers []string

	rows, err := db.Query("SELECT question, answer FROM history ORDER BY time DESC")

	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		var question string
		var answer string
		rows.Scan(&question, &answer)
		questions = append(questions, question)
		answers = append(answers, answer)
	}

	return questions, answers
}

func Read_History(db *sql.DB, question string) string { // Marked for deletion
	var answer string
	db.QueryRow("SELECT answer FROM history WHERE question = ?", question).Scan(&answer)
	return answer
}

func Delete_History(db *sql.DB, question string) { // Marked for deletion
	_, err := db.Query("DELETE FROM history WHERE question = ?", question)

	if err != nil {
		panic(err.Error())
	}
}

// CRUD operations for questions table

func Create_Question(db *sql.DB, question string, answer string) {
	_, err := db.Query("INSERT IGNORE INTO questions (question, answer) VALUES (?, ?)", question, answer)

	if err != nil {
		panic(err.Error())
	}
}

func Read_Question(db *sql.DB, question string) string {
	var answer string
	db.QueryRow("SELECT answer FROM questions WHERE question = ?", question).Scan(&answer)
	return answer
}

func Read_All_Questions(db *sql.DB) []string {
	var questions []string
	rows, err := db.Query("SELECT question FROM questions")

	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		var question string
		rows.Scan(&question)
		questions = append(questions, question)
	}

	return questions
}

func Update_Question(db *sql.DB, question string, answer string) { // Marked for deletion
	Delete_Question(db, question)
	Create_Question(db, question, answer)
}

func Update_Answer(db *sql.DB, question string, answer string) {
	_, err := db.Query("UPDATE questions SET answer = ? WHERE question = ?", answer, question)

	if err != nil {
		panic(err.Error())
	}
}

func Delete_Question(db *sql.DB, question string) {
	_, err := db.Query("DELETE FROM questions WHERE question = ?", question)

	if err != nil {
		panic(err.Error())
	}
}
