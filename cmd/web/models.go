package main

type Word struct {
	ID      string `json:"id"`
	English string `json:"english"`
	Polish  string `json:"polish"`
}

type errorMessage struct {
	Message string
}
