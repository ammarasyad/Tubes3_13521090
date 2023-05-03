package sql_connection

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func create_history_table(db *sql.DB) {
	_, err := db.Query(`CREATE TABLE IF NOT EXISTS history (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		question VARCHAR(255) NOT NULL,
		answer VARCHAR(255) NOT NULL,
		time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)`)

	if err != nil {
		panic(err.Error())
	}
}

// CRD operations for history table (no update)

func create_history(db *sql.DB, question string, answer string) {
	_, err := db.Query("INSERT IGNORE INTO history (question, answer) VALUES (?, ?)", question, answer)

	if err != nil {
		panic(err.Error())
	}
}

func read_history(db *sql.DB, question string) string {
	var answer string
	db.QueryRow("SELECT answer FROM history WHERE question = ?", question).Scan(&answer)
	return answer
}

func delete_history(db *sql.DB, question string) {
	_, err := db.Query("DELETE FROM history WHERE question = ?", question)

	if err != nil {
		panic(err.Error())
	}
}

func create_questions_table(db *sql.DB) {
	_, err := db.Query(`CREATE TABLE IF NOT EXISTS questions (
		question VARCHAR(255) NOT NULL PRIMARY KEY,
		answer VARCHAR(255) NOT NULL)`)

	if err != nil {
		panic(err.Error())
	}
}

// CRUD operations for questions table

func create_question(db *sql.DB, question string, answer string) {
	_, err := db.Query("INSERT IGNORE INTO questions (question, answer) VALUES (?, ?)", question, answer)

	if err != nil {
		panic(err.Error())
	}
}

func read_question(db *sql.DB, question string) string {
	var answer string
	db.QueryRow("SELECT answer FROM questions WHERE question = ?", question).Scan(&answer)
	return answer
}

func update_question(db *sql.DB, question string, answer string) {
	_, err := db.Query("UPDATE questions SET answer = ? WHERE question = ?", answer, question)

	if err != nil {
		panic(err.Error())
	}
}

func delete_question(db *sql.DB, question string) {
	_, err := db.Query("DELETE FROM questions WHERE question = ?", question)

	if err != nil {
		panic(err.Error())
	}
}
