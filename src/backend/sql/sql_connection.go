package sql_connection

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Create_Database(conn *sql.Conn, ctx context.Context) {
	_, err := conn.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS GyrosPallas")

	if err != nil {
		panic(err.Error())
	}

	if err != nil {
		panic(err.Error())
	}

	_, err = conn.ExecContext(ctx, "USE GyrosPallas")

	if err != nil {
		panic(err.Error())
	}

	_, err = conn.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS history (
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		question VARCHAR(255) NOT NULL,
		answer VARCHAR(255) NOT NULL,
		time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)`)

	if err != nil {
		panic(err.Error())
	}

	_, err = conn.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS questions (
		question VARCHAR(255) NOT NULL PRIMARY KEY,
		answer VARCHAR(255) NOT NULL)`)

	if err != nil {
		panic(err.Error())
	}
}

// CRD operations for history table (no update)

func Create_History(conn *sql.Conn, ctx context.Context, question string, answer string) {
	_, err := conn.QueryContext(ctx, "INSERT IGNORE INTO history (question, answer) VALUES (?, ?)", question, answer)

	if err != nil {
		panic(err.Error())
	}
}

func Read_History(conn *sql.Conn, ctx context.Context, question string) string {
	var answer string
	conn.QueryRowContext(ctx, "SELECT answer FROM history WHERE question = ?", question).Scan(&answer)
	return answer
}

func Delete_History(conn *sql.Conn, ctx context.Context, question string) {
	_, err := conn.QueryContext(ctx, "DELETE FROM history WHERE question = ?", question)

	if err != nil {
		panic(err.Error())
	}
}

// CRUD operations for questions table

func Create_Question(conn *sql.Conn, ctx context.Context, question string, answer string) {
	_, err := conn.QueryContext(ctx, "INSERT IGNORE INTO questions (question, answer) VALUES (?, ?)", question, answer)

	if err != nil {
		panic(err.Error())
	}
}

func Read_Question(conn *sql.Conn, ctx context.Context, question string) string {
	var answer string
	conn.QueryRowContext(ctx, "SELECT answer FROM questions WHERE question = ?", question).Scan(&answer)
	return answer
}

func Read_All_Questions(conn *sql.Conn, ctx context.Context) []string {
	var questions []string
	rows, err := conn.QueryContext(ctx, "SELECT question FROM questions")

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

func Update_Question(conn *sql.Conn, ctx context.Context, question string, answer string) {
	_, err := conn.QueryContext(ctx, "UPDATE questions SET answer = ? WHERE question = ?", answer, question)

	if err != nil {
		panic(err.Error())
	}
}

func Update_Answer(conn *sql.Conn, ctx context.Context, question string, answer string) {
	_, err := conn.QueryContext(ctx, "UPDATE questions SET answer = ? WHERE question = ?", answer, question)

	if err != nil {
		panic(err.Error())
	}
}

func Delete_Question(conn *sql.Conn, ctx context.Context, question string) {
	_, err := conn.QueryContext(ctx, "DELETE FROM questions WHERE question = ?", question)

	if err != nil {
		panic(err.Error())
	}
}
