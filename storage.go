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
	return s.CreateTables()
}

func (s *PostgresStore) CreateTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS language (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) UNIQUE,
		usage NUMERIC(5, 2)
	  );

	  CREATE TABLE IF NOT EXISTS code_report (
		id SERIAL PRIMARY KEY,
		request INTEGER,
		language_id INTEGER REFERENCES language(id),
		score INTEGER,
		usage NUMERIC(5, 2),
		created_at TIMESTAMPTZ DEFAULT NOW()
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

func (s *PostgresStore) CreateCodeReport(cr *CodeReport) error {
	query := `
	INSERT INTO code_report (request, language_id, score, usage, created_at)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.Exec(query, cr.Request, cr.Language_ID, cr.Score, cr.Use_Percentage, cr.Created_At)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) CreateLanguage(l *Language) error {
	query := `
	INSERT INTO language (name)
	VALUES ($1)`
	_, err := s.db.Exec(query, l.Name)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStore) getLanguageId(name string) (int, error) {
	var id int

	err := s.db.QueryRow("SELECT id FROM language WHERE name = $1", name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return id, nil
}

func (s *PostgresStore) UpdateLanguageUsage(l *Language) error {
	// Begin the transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// Prepare the update statement
	stmt, err := tx.Prepare("UPDATE language SET usage = $1 WHERE name = $2")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	// Execute the update statement
	_, err = stmt.Exec(l.Usage, l.Name)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (s *PostgresStore) getLastRequest() (int, error) {
	var lastGroupID int
	err := s.db.QueryRow("SELECT COALESCE(MAX(request), 0) FROM code_report").Scan(&lastGroupID)
	if err != nil {
		lastGroupID = -1
		return 0, err
	}
	return lastGroupID, nil
}
