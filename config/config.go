package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := "user=postgres password=postgres dbname=db_go sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS items (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		body TEXT NOT NULL
	);
	`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database connected")
}
