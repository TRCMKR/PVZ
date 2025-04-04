package config

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

// InitEnv ...
func InitEnv(envFile string) {
	err := godotenv.Overload(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file", envFile)
	}
}

// Config ...
type Config struct {
	host        string
	port        string
	username    string
	password    string
	dbname      string
	WorkerCount int
	BatchSize   int
	Timeout     time.Duration
}

// NewConfig ...
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
		host:        host,
		port:        port,
		username:    username,
		password:    password,
		dbname:      dbname,
		WorkerCount: 2,
		BatchSize:   5,
		Timeout:     2 * time.Second,
	}
}

func (c *Config) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.host, c.port, c.username, c.password, c.dbname)
}

// Host ...
func (c *Config) Host() string {
	return c.host
}

// Port ...
func (c *Config) Port() string {
	return c.port
}

// Username ...
func (c *Config) Username() string {
	return c.username
}

// Password ...
func (c *Config) Password() string {
	return c.password
}

// DBName ...
func (c *Config) DBName() string {
	return c.dbname
}

// GetRootDir ...
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

// ReadFirstFileWord ...
func ReadFirstFileWord(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	if !scanner.Scan() {
		return "", errors.New("file is empty or has no words")
	}
	firstWord := scanner.Text()

	if scanner.Scan() {
		return "", errors.New("file contains more than one word")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return firstWord, nil
}
