package types

import "time"

type Language struct {
  ID int `json:id`
  Name  string `json:"name"` 
  Usage float64 `json:"usage"`
}

type CodeReport struct {
	Request_ID     int
	Language_ID    int
	Score          int
	Use_Percentage float64
	Created_At     time.Time
}


