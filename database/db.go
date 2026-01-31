package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error

	// Open or create DB file
	DB, err = sql.Open("sqlite3", "./drone-system.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Test DB connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	log.Println("Database connected")

	createTables()
}

func createTables() {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		type TEXT
	);`

	dronesTable := `
	CREATE TABLE IF NOT EXISTS drones (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
		status TEXT,
		lat REAL,
		lng REAL,
		current_order_id INTEGER
	);`

	ordersTable := `
	CREATE TABLE IF NOT EXISTS orders (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
		origin TEXT,
		destination TEXT,
		status TEXT,
		assigned_drone_id INTEGER,
		current_lat REAL,
		current_lng REAL,
		created_by INTEGER
	);`

	_, err := DB.Exec(usersTable)
	if err != nil {
		log.Fatal("Failed creating users table:", err)
	}

	_, err = DB.Exec(dronesTable)
	if err != nil {
		log.Fatal("Failed creating drones table:", err)
	}

	_, err = DB.Exec(ordersTable)
	if err != nil {
		log.Fatal("Failed creating orders table:", err)
	}

	log.Println("Tables created (if not exist)")
}
