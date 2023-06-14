package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Storage interface {
	createAccount(*User) error
}

type PostgresStore struct {
	db *sql.DB
}

func (s *PostgresStore) Init() error {
	return s.CreateLanguageTable()
}

func (s *PostgresStore) CreateLanguageTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS languages (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) UNIQUE
	  );

	CREATE TABLE IF NOT EXISTS code_report (
		id SERIAL PRIMARY KEY,
		request INTEGER,
		language_id INTEGER REFERENCES languages(id),
		score INTEGER,
		created_at TIMESTAMPTZ DEFAULT NOW(),
		percentage NUMERIC(5, 5)
	);
`
	_, err := s.db.Exec(query)
	return err
}

func NewPostgresStore() (*PostgresStore, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	user := os.Getenv("USER_DB")
	password := os.Getenv("PASSWORD")
	database := os.Getenv("DATABASE")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, database)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

// func (s *PostgresStore) CreateCodeReport
