package main

import "time"

type Language struct {
	Name string
}

type CodeReport struct {
	Request     int
	Language_id int
	Score       int
	Percentage  float64
	Created_At  time.Time
}
