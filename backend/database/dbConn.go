package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"skillswap/backend/utils"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Init() {
	var err error
	db, err = sql.Open("mysql", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	utils.DebugPrint("Database connected")
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	utils.DebugPrint("Database ping successful")
	res, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err.Error())
		if strings.Contains(err.Error(), "doesn't exist"){
			Migrate(db)
		}
	}
	res.Close()
}

func GetDatabase() (*sql.DB, error) {
	return db, nil
}

func Execute(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}

func QueryRow(query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}

func Close() error {
	return db.Close()
}

func Debug(query string, args ...interface{}) {
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
