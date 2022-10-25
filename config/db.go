package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var db *gorm.DB

func InitDB() {
	var err error
	db, err = gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect to the database: \n", err)
	}
}

func GetDB() *gorm.DB {
	return db
}

func CloseDB() {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln("Failed to close connection to the database: \n", err)
	}

	sqlDB.Close()
}
