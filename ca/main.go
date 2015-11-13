package ca

import (
	"flag"
	"log"
)

func main() {
	// showing the source file and line number where the log statement comes from
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()

	app, err := NewApp()
	if err != nil {
		log.Fatal(err)
	}

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
