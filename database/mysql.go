package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectToMySQLDB connects to the database
func ConnectToMySQLDB() (*sql.DB, error) {
	dsn := "username:password@tcp(127.0.0.1:3306)/dbname" //This string can be obtained from the env file
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

// ExecuteMySQLQuery executes a query and returns the result
func ExecuteMySQLQuery(db *sql.DB, query string, args ...any) (*sql.Rows, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	return rows, nil
}

// ExecuteMySQLNonQuery execute operations like INSERT, UPDATE, DELETE
func ExecuteMySQLNonQuery(db *sql.DB, query string, args ...any) (sql.Result, error) {
	result, err := db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute non-query: %v", err)
	}
	return result, nil
}
