package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func InitEnv(envFile string) {
	err := godotenv.Overload(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file", envFile)
	}
}

type Config struct {
	host     string
	port     string
	username string
	password string
	dbname   string
}

func NewConfig() *Config {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || username == "" || password == "" || dbname == "" {
		log.Fatal("Database configuration missing: one or more required fields are empty.")
	}

	return &Config{
		host:     host,
		port:     port,
		username: username,
		password: password,
		dbname:   dbname,
	}
}

func (c *Config) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.host, c.port, c.username, c.password, c.dbname)
}

func (c *Config) Host() string {
	return c.host
}

func (c *Config) Port() string {
	return c.port
}

func (c *Config) Username() string {
	return c.username
}

func (c *Config) Password() string {
	return c.password
}

func (c *Config) DBName() string {
	return c.dbname
}

func GetRootDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err = os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}
