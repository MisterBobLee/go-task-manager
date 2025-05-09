package models

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type User struct {
    gorm.Model
    Username string `gorm:"unique"`
    Password string
}

func ConnectDatabase() {
    dsn := os.Getenv("DB_URL")
    if dsn == "" {
        log.Fatal("Environment variable DB_URL not set")
    }
    
    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    fmt.Println("Connected to database")

    database.AutoMigrate(&User{}, &Task{})
    DB = database
}
