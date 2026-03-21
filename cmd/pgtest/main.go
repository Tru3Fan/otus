package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	dsn := "host=localhost port=5432 user=admin password=secret dbname=otus sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println("Open error:", err)
		return
	}
	if err := db.Ping(); err != nil {
		fmt.Println("Ping error:", err)
		return
	}
	fmt.Println("Connected!")
}
