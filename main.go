package main

import (
	"flag"
	"log"

	"github.com/tslilyai/SLYcoin/app"
)

func main() {
	log.Println("Welcome to SLYcoin!")
	flag.Parse()

	app, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
