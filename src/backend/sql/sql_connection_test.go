package sql_connection

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestTable(t *testing.T) {
	data_source := "data source here"
	db, err := sql.Open("mysql", data_source)

	if err != nil {
		panic(err.Error())
	}

	ctx := context.Background()
	conn, err := db.Conn(ctx)

	if err != nil {
		panic(err.Error())
	}

	Create_Database(conn, ctx)

	Create_History(conn, ctx, "apa", "itu")
	Create_Question(conn, ctx, "apa", "itu")

	if Read_History(conn, ctx, "apa") != "itu" {
		t.Error("read_history() failed")
	}

	Update_Question(conn, ctx, "apa", "itu apa")

	if Read_Question(conn, ctx, "apa") != "itu apa" {
		t.Error("read_question() failed")
	}

	Delete_History(conn, ctx, "apa")
	Delete_Question(conn, ctx, "apa")

	defer db.Close()
}

func TestFile(t *testing.T) {
	filename := "stima.sql"

	fileInfo, err := os.Stat(filename)

	if err != nil {
		panic(err.Error())
	}

	if fileInfo.Size() == 0 {
		t.Error("File is empty")
	}

	fmt.Println(fileInfo.Size())
	fmt.Println(fileInfo.Name())
	fmt.Println(fileInfo.IsDir())
	fmt.Println(fileInfo.ModTime())
	fmt.Println(fileInfo.Mode())
	fmt.Println(fileInfo.Sys())
}
