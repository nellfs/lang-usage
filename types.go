package main

import "time"

type Language struct {
	ID   int
	Name string
}

type CodeReport struct {
	ID          int
	Request     int
	Language_id int
	score       int
	created_at  time.Time
	percentage  float64
}
