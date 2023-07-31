package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/thisguycodes/dyno/sequences"
	"github.com/thisguycodes/dyno/structure"
)

type Flags struct {
	DBName        string
	MaxDistance   int
	ListenAddress string
}

func parseFlags() Flags {
	flags := Flags{}

	flag.StringVar(&flags.DBName, "DBName", "db.sqlite3", "Name of the sqlite3 database file to use")
	flag.IntVar(&flags.MaxDistance, "MaxDistance", 1, "The max hamming distance to query for")
	flag.StringVar(&flags.ListenAddress, "address", ":8080", "Address to listen for requests on")
	flag.Parse()

	return flags
}

func main() {
	ctx := context.Background()

	flags := parseFlags()

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared", flags.DBName))
	if err != nil {
		log.Fatalf("Couldn't open DB: %q", err)
	}

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Couldn't ping DB: %q", err)
	}

	if err := structure.InitTables(ctx, db); err != nil {
		log.Fatalf("Could not create tables: %q", err)
	}

	uploader := sequences.Uploader{
		DB: db,
	}
	querier := sequences.HammingQuery{
		DB:          db,
		MaxDistance: flags.MaxDistance,
	}

	http.Handle("/upload", uploader)
	http.Handle("/hamming_matches", querier)

	if err := http.ListenAndServe(flags.ListenAddress, nil); err != nil {
		log.Printf("Error from http.ListenAndServe: %q", err)
	}
}
