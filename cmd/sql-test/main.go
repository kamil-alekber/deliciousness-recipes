package main

import (
	_ "embed"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

//go:embed sql/app/schema.sql
var ddl string

func run() error {
	fmt.Println("setup cron job for sql workflows")

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
