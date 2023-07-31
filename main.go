package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/thisguycodes/dyno/structure"
)

type Flags struct {
	DBName string
}

func parseFlags() Flags {
	flags := Flags{}

	flag.StringVar(&flags.DBName, "DBName", "db.sqlite3", "Name of the sqlite3 database file to use")
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
}
