package main

import (
	"log"

	"brianhang.me/facegraph/internal/db"
	"brianhang.me/facegraph/internal/routes"
	"brianhang.me/facegraph/internal/user"
	"gorm.io/gorm"
)

func main() {
	db, err := db.Init()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}

	if err = setupModels(db); err != nil {
		log.Fatalf("failed to set up models: %v", err)
	}

	err = routes.Init()
	if err != nil {
		log.Fatalf("Failed to initialize routes: %v", err)
	}
}

func setupModels(db *gorm.DB) error {
	return db.AutoMigrate(
		user.User{},
	)
}
