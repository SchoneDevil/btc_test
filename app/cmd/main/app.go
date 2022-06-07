package main

import (
	"log"

	"app/internal/app"
)

func main() {
	a, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Preload Running")
	a.Preload()
	log.Println("Running Application")
	a.Run()
}
