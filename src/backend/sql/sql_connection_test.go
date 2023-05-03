package sql_connection

import (
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestHistoryTable(t *testing.T) {
	db, err := sql.Open("mysql", "root:Archr181003.@/stima")

	if err != nil {
		panic(err.Error())
	}

	create_history_table(db)
	create_questions_table(db)

	create_history(db, "apa", "itu")
	create_question(db, "apa", "itu")

	if read_history(db, "apa") != "itu" {
		t.Error("read_history() failed")
	}

	update_question(db, "apa", "itu apa")

	if read_question(db, "apa") != "itu apa" {
		t.Error("read_question() failed")
	}

	delete_history(db, "apa")
	delete_question(db, "apa")

	defer db.Close()
}
