/*
** D&GINE Project, 2026
** API
** File description:
** database/database.go
 */

package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var Db *sql.DB

// CloseRows closes rows and logs the error if closing fails.
// Intended to be used with defer; iteration errors must still be checked with rows.Err().
func CloseRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Printf("database: failed to close rows: %v", err)
	}
}

func ConnectDatabase() {
	host := os.Getenv("DB_HOST")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PWD")
	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, pass)
	db, errSql := sql.Open("postgres", psqlSetup)
	err := db.Ping()
	if errSql != nil {
		fmt.Println("There is an error while connecting to the database ", errSql)
		panic(errSql)
	} else if err != nil {
		panic(err)
	} else {
		Db = db
		fmt.Println("Successfully connected to database!")
	}
}
