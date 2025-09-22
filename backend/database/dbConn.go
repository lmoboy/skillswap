package database


import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func getDatabase() (*sql.DB, error) {
	return db, nil
}

func execute(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}

func query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func queryRow(query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}

func close() error {
	return db.Close()
}

func debug(query string, args ...interface{}) {
	fmt.Printf("Query: %s\nArgs: %v\n", query, args)
}

func logError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func logQueryError(err error) {
	if err != nil {
		log.Println("Query error:", err)
	}
}

func logQueryResult(result sql.Result, err error) {
	if err != nil {
		log.Println("Query result:", result, err)
	}
}

func logQueryRow(row *sql.Row) {
	if row == nil {
		log.Println("Query row: nil")
	} else {
		log.Println("Query row:", row)
	}
}
