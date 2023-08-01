package main

import (
	"time"
)

type Language struct {
  Name  string 
  Usage float64 
}

type CodeReport struct {
  Request_ID      int 
	Language_ID    int 
	Score          int
	Use_Percentage float64
	Created_At     time.Time
}

