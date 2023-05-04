package main

import (
	"database/sql"
	algorithm "src/backend/src/backend/algorithm"
	sql_connection "src/backend/src/backend/sql"
)

func main() {
	data_source := "root:password@/"
	db, err := sql.Open("mysql", data_source)

	if err != nil {
		panic(err.Error())
	}

	sql_connection.Create_Database(db)
	println(algorithm.ProcessQuestion(db, "tambahkan pertanyaan papope dengan jawaban poppop", true))
}
