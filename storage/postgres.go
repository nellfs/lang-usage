package storage;

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nellfs/lang-usage/types"
  	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	DB *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
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
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStorage{
		DB: db,
	}, nil
}

func (ps *PostgresStorage) Init() error {
	query := `
	CREATE TABLE IF NOT EXISTS languages (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) UNIQUE,
		usage NUMERIC(5, 2)
	  );

	  CREATE TABLE IF NOT EXISTS code_reports (
		id SERIAL PRIMARY KEY,
		request INTEGER,
		language_id INTEGER REFERENCES languages(id),
		score INTEGER,
		usage NUMERIC(5, 2),
		created_at TIMESTAMPTZ DEFAULT NOW()
	);

`
	_, err := ps.DB.Exec(query)
	return err
}

func (ps *PostgresStorage) CreateLanguage(l *types.Language) error {
	query := `
	INSERT INTO languages (name)
	VALUES ($1)`
	_, err := ps.DB.Exec(query, l.Name)
	if err != nil {
		return err
	}
	return nil
}

//Return list with all languages
func (ps *PostgresStorage) GetLanguage(language *string) ([]*type.Language, error) {
  // rows, err := ps.DB.Query("select * from languages")
  // if err != nil {
  //   return nil, err
  // }
  // rows.Next()
}


func (ps *PostgresStorage) GetLanguageIDByName(name string) (int, error) {
	var id int

	err := ps.DB.QueryRow("SELECT id FROM languages WHERE name = $1", name).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return id, nil
}

func (ps *PostgresStorage) CreateCodeReport(cr *types.CodeReport) error {
	query := `
	INSERT INTO code_reports (request, language_id, score, usage, created_at)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := ps.DB.Exec(query, cr.Request_ID, cr.Language_ID, cr.Score, cr.Use_Percentage, cr.Created_At)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostgresStorage) GetCodeReport(number int) (*types.CodeReport, error) {
	return nil, nil
}

func (ps *PostgresStorage) GetLastRequestID() (int, error) {
	var lastGroupID int
	err := ps.DB.QueryRow("SELECT COALESCE(MAX(request), 0) FROM code_reports").Scan(&lastGroupID)
	if err != nil {
		lastGroupID = -1
		return 0, err
	}
	return lastGroupID, nil
}
