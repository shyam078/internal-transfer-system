package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	var err error

	// Load environment variables from .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read database connection parameters from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("Database environment variables are not set properly")
	}

	// Build connection string dynamically
	connStr := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable"

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}
	defer DB.Close()

	// Check if the target database exists
	var exists bool
	err = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = 'internal_transfer')").Scan(&exists)
	if err != nil {
		log.Fatal("Error checking database existence:", err)
	}

	// Create the database if it doesn't exist
	if !exists {
		_, err = DB.Exec("CREATE DATABASE internal_transfer")
		if err != nil {
			log.Fatal("Error creating database:", err)
		}
		log.Println("Database 'internal_transfer' created.")
	}

	// Now connect to the target database
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to target DB:", err)
	}

	createTables()
}

func createTables() {
	_, err := DB.Exec(`
        CREATE TABLE IF NOT EXISTS accounts (
            account_id BIGINT PRIMARY KEY,
            balance NUMERIC NOT NULL
        );
        CREATE TABLE IF NOT EXISTS transactions (
            id SERIAL PRIMARY KEY,
            source_account_id BIGINT,
            destination_account_id BIGINT,
            amount NUMERIC,
            created_at TIMESTAMP DEFAULT now()
        );
    `)
	if err != nil {
		log.Fatal("Error creating tables:", err)
	}
}
