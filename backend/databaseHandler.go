package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// addValues adds the given values to the specified table in the database.
//
// The values will be inserted in the order specified by the fields argument.
// The function will return if there is an error with the database connection
// or if there is an error with the SQL statement.
//
// Usage:
//
// addValues("table", []string{"field1", "field2"}, [][]string{{"value1", "value2"}})

func getDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/skillswap")
	return db, err
}

func addValues(table string, fields []string, values [][]string) error {
	parsedFields := strings.Join(fields, ", ")
	fieldsToFill := strings.Repeat("?,", len(fields))
	fieldsToFill = fieldsToFill[:len(fieldsToFill)-1]

	db, err := getDatabase()
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO " + table + " (" + parsedFields + ") VALUES (" + fieldsToFill + ")")

	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()

	for _, row := range values {
		// Convert the []string to []interface{}
		args := make([]interface{}, len(row))
		for i, v := range row {
			args[i] = v
		}

		// Execute the prepared statement with the variadic arguments
		_, err = stmt.Exec(args...)
		if err != nil {
			fmt.Printf("Error executing statement for row %v: %v\n", row, err)
			return err
		}
	}
	return nil
}
func removeValues(table string, conditions map[string]string) error {
	var conds []string
	var args []interface{}
	for k, v := range conditions {
		conds = append(conds, fmt.Sprintf("%s = ?", k))
		args = append(args, v)
	}
	whereClause := strings.Join(conds, " AND ")

	db, err := getDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM " + table + " WHERE " + whereClause)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}

	return nil
}
func updateValues(table string, updates map[string]string, conditions map[string]string) error {
	var setClauses []string
	var args []interface{}
	for k, v := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", k))
		args = append(args, v)
	}
	setClause := strings.Join(setClauses, ", ")

	var condClauses []string
	for k, v := range conditions {
		condClauses = append(condClauses, fmt.Sprintf("%s = ?", k))
		args = append(args, v)
	}
	whereClause := strings.Join(condClauses, " AND ")

	db, err := getDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE " + table + " SET " + setClause + " WHERE " + whereClause)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}

	return nil
}
func findValues(table string, values []string, conditions map[string]string) ([]map[string]string, error) {
	// findValues finds values in the database that match the given conditions
	if table == "" {
		return nil, fmt.Errorf("table name is empty")
	}
	if len(conditions) == 0 {
		return nil, fmt.Errorf("conditions map is empty")
	}

	var condClauses []string
	var args []interface{}

	for k, v := range conditions {
		condClauses = append(condClauses, fmt.Sprintf("%s = ?", k))
		args = append(args, v)
	}
	whereClause := strings.Join(condClauses, " AND ")
	selectClause := "*"
	if len(values) > 0 {
		selectClause = strings.Join(values, ", ")
	}

	db, err := getDatabase()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT " + selectClause + " FROM " + table + " WHERE " + whereClause)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	rvalues := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(rvalues))
	for i := range rvalues {
		scanArgs[i] = &rvalues[i]
	}

	var results []map[string]string
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		row := make(map[string]string, len(columns))
		for i, col := range columns {
			row[col] = string(rvalues[i])
		}
		results = append(results, row)
	}

	return results, nil
}

// func main() {
// 	err := addValues("users", []string{"username", "email", "password_hash"}, [][]string{{"username1", "email1@email.email", "password_hash1"}, {"username2", "email2@email.email", "password_hash2"}})
// 	if err != nil {
// 		fmt.Println("Error adding users:", err)
// 	} else {
// 		fmt.Println("Users added successfully.")
// 	}
// 	err = removeValues("users", map[string]string{"username": "username1"})
// 	if err != nil {
// 		fmt.Println("Error removing user:", err)
// 	} else {
// 		fmt.Println("User 'username1' removed successfully.")
// 	}

// 	err = updateValues("users", map[string]string{"email": "new_email2@email.email"}, map[string]string{"username": "username2"})
// 	if err != nil {
// 		fmt.Println("Error updating user:", err)
// 	} else {
// 		fmt.Println("User 'username2' updated successfully.")
// 	}
// }
