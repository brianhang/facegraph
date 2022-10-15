package main

import (
	"log"

	"brianhang.me/facegraph/internal/db"
	"brianhang.me/facegraph/internal/routes"
)

func main() {
	err := db.Init()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}

	err = routes.Init()
	if err != nil {
		log.Fatalf("Failed to initialize routes: %v", err)
	}
}
