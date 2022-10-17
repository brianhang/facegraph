package db

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbInstance *gorm.DB

func Init() (*gorm.DB, error) {
	dbFileName := os.Getenv("DB_FILE_NAME")
	if dbFileName == "" {
		dbFileName = "app.db"
	}

	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})

	log.Printf("Using sqlite with main database file at %s\n", dbFileName)

	dbInstance = db

	return db, err
}

func Get() *gorm.DB {
	return dbInstance
}
