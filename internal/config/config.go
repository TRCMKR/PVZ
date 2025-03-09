package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

//type Config struct {
//	DBHost     string
//	DBPort     string
//	DBUsername string
//	DBPassword string
//	DBName     string
//}
//
//func NewConfig() *Config {
//	return &Config{
//		DBHost:     os.Getenv("DB_HOST"),
//		DBPort:     os.Getenv("DB_PORT"),
//		DBUsername: os.Getenv("DB_USER"),
//		DBPassword: os.Getenv("DB_PASSWORD"),
//		DBName:     os.Getenv("DB_NAME"),
//	}
//}

func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func GetDBHost() string {
	return os.Getenv("DB_HOST")
}

func GetDBPort() string {
	return os.Getenv("DB_PORT")
}

func GetDBUsername() string {
	return os.Getenv("DB_USERNAME")
}

func GetDBPassword() string {
	return os.Getenv("DB_PASSWORD")
}

func GetDBName() string {
	return os.Getenv("DB_NAME")
}
