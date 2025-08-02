package config

import (
	"appGO/model"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var JwtKey []byte

func LoadEnv() {
	err := godotenv.Load()             
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	JwtKey = []byte(os.Getenv("JWT_SECRET"))
}

func InitDB() {
	dsn := "host=localhost port=5432 user=farhanyousuf dbname=goapp sslmode=disable"
	var err error

	
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	

	log.Println("✅ Connected to database using GORM")


    // Auto migrate multiple models
    err = DB.AutoMigrate(
        &model.User{},
        &model.Category{},
        &model.Item{},
        // &model.Order{},
        // // Add more models here as needed
    )
    if err != nil {
        log.Fatalf("AutoMigrate failed: %v", err)
    }
}
