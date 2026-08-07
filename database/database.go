package database

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Connect() {
	var err error

	DB, err = sql.Open("sqlite", "./database/myjob.db")
	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Connected to SQLite")

	createTables()
}

func createTables() {
	schema, err := os.ReadFile("database/schema.sql")
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(string(schema))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Database tables created successfully")
}
