package main

import (
	"flag"
	"log"
	"time"

	"github.com/diother/hintermann-stripe-docgen/internal/app"
)

func main() {
	year := flag.Int("year", time.Now().Year(), "Year for monthly report")
	month := flag.Int("month", int(time.Now().Month()), "Month for monthly report")
	flag.Parse()

	if err := app.Run(*year, *month); err != nil {
		log.Fatal(err)
	}
}
