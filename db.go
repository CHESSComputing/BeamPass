package main

import (
	"database/sql"
	"fmt"
	"log"
)

var db *sql.DB

func initDB(dbURI string) {
	var err error
	db, err = sql.Open("mysql", dbURI)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	// Ping the database to ensure connection is established
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	fmt.Println("Successfully connected to MySQL!")
}
