package main

import "database/sql"


type Storage interface {
	createAccount(*User) error
}

type PostgresStore struct {
	db *sql.DB
i