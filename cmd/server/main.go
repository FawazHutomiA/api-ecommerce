package main

import (
	"example/config"
	"example/internal/router"
	"log"
)

func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	router.Router(db)
}
