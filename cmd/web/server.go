package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var cfg struct {
	addr *string
	dsn  *string
	wait *int
}

type DictionaryDatabase struct {
	db *sql.DB
}

var lidiDB DictionaryDatabase

func main() {
	// parametry serwera z linii komend, domyślnie port 8080 a dsn ze zmiennej środowiskowej
	cfg.addr = flag.String("addr", ":8080", "HTTP network address")
	cfg.wait = flag.Int("wait", 120, "time to wait for server (in sec)")
	cfg.dsn = flag.String("dsn", os.Getenv("LIDI_SERVER_SECRET"), "MySQL data source name")
	flag.Parse()

	// połączenie z bazą danych
	err := lidiDB.openDB(*cfg.dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer lidiDB.db.Close()

	// router, definicje endpointów
	router := RegisterRoutes()

	// start serwera
	fmt.Printf("Server started, port %s", *cfg.addr)
	http.ListenAndServe(*cfg.addr, router)
	log.Fatal(err)
}
