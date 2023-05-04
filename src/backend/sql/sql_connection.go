package sql_connection

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Create_Database(db *sql.DB) {
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS GyrosPallas")

	if err != nil {
		panic(err.Error())
	}

	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS GyrosPallas.history (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		question VARCHAR(255) NOT NULL,
		answer VARCHAR(255) NOT NULL,
		time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)`)

	if err != nil {
		panic(err.Error())
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS GyrosPallas.questions (
		question VARCHAR(255) NOT NULL PRIMARY KEY,
		answer VARCHAR(255) NOT NULL)`)

	if err != nil {
		panic(err.Error())
	}
}

// CRD operations for history table (no update)

func Create_History(db *sql.DB, question string, answer string) {
	_, err := db.Exec("INSERT IGNORE INTO GyrosPallas.history (question, answer) VALUES (?, ?)", question, answer)

	if err != nil {
		panic(err.Error())
	}
}

func Read_History(db *sql.DB, question string) string {
	var answer string
	db.QueryRow("SELECT answer FROM GyrosPallas.history WHERE question = ?", question).Scan(&answer)
	return answer
}

func Delete_History(db *sql.DB, question string) {
	_, err := db.Query("DELETE FROM GyrosPallas.history WHERE question = ?", question)

	if err != nil {
		panic(err.Error())
	}
}

// CRUD operations for questions table

func Create_Question(db *sql.DB, question string, answer string) {
	_, err := db.Query("INSERT IGNORE INTO GyrosPallas.questions (question, answer) VALUES (?, ?)", question, answer)

	if err != nil {
		panic(err.Error())
	}
}

func Read_Question(db *sql.DB, question string) string {
	var answer string
	db.QueryRow("SELECT answer FROM GyrosPallas.questions WHERE question = ?", question).Scan(&answer)
	return answer
}

func Read_All_Questions(db *sql.DB) []string {
	var questions []string
	rows, err := db.Query("SELECT question FROM GyrosPallas.questions")

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

func Update_Question(db *sql.DB, question string, answer string) {
	Delete_Question(db, question)
	Create_Question(db, question, answer)
}

func Update_Answer(db *sql.DB, question string, answer string) {
	_, err := db.Query("UPDATE GyrosPallas.questions SET answer = ? WHERE question = ?", answer, question)

	if err != nil {
		panic(err.Error())
	}
}

func Delete_Question(db *sql.DB, question string) {
	_, err := db.Query("DELETE FROM GyrosPallas.questions WHERE question = ?", question)

	if err != nil {
		panic(err.Error())
	}
}
