package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var cfg struct {
	addr *string
	dsn  *string
}

var db *sql.DB
var err error

func main() {
	// parametry serwera z linii komend
	cfg.addr = flag.String("addr", ":8080", "HTTP network address")
	cfg.dsn = flag.String("dsn", "web:pass@/dictionary", "MySQL data source name")
	flag.Parse()

	// połączenie z bazą danych
	db, err = openDB(*cfg.dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// router, definicje endpointów
	router := RegisterRoutes()

	// start serwera
	fmt.Printf("Server started, port %s", *cfg.addr)
	http.ListenAndServe(*cfg.addr, router)
	log.Fatal(err)
}
