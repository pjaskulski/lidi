package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Word struct {
	ID      string `json:"id"`
	English string `json:"english"`
	Polish  string `json:"polish"`
}

var db *sql.DB
var err error

func openDB(dsn string) (*sql.DB, error) {
	hDb, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = hDb.Ping(); err != nil {
		return nil, err
	}
	return hDb, nil
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address")
	dsn := flag.String("dsn", "web:pass@/dictionary", "MySQL data source name")
	flag.Parse()

	db, err = openDB(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/api/en/{word}", getPolish).Methods("GET")
	router.HandleFunc("/api/add", createWord).Methods("POST")
	router.HandleFunc("/api/pl/{word}", getEnglish).Methods("GET")
	//router.HandleFunc("/api/update", updateWord).Methods("PUT")
	//router.HandleFunc("/api/delete", deleteWord).Methods("DELETE")

	fmt.Printf("Start serwera, port %s", *addr)
	http.ListenAndServe(*addr, router)
	log.Fatal(err)
}
