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

// InitEnv inits env file from path
func InitEnv(envFile string) {
	err := godotenv.Overload(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file", envFile)
	}
}

// Config is a structure that contains all configuration parameters
type Config struct {
	host          string
	port          string
	username      string
	password      string
	dbname        string
	kafka_port    string
	kafka_ui_port string
	WorkerCount   int
	BatchSize     int
	Timeout       time.Duration
}

// NewConfig creates instance of Config
func NewConfig() *Config {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	kafka_port := os.Getenv("KAFKA_PORT")
	kafka_ui_port := os.Getenv("KAFKA_UI_PORT")

	if host == "" || port == "" || username == "" || password == "" || dbname == "" ||
		kafka_port == "" || kafka_ui_port == "" {
		log.Fatal("Database configuration missing: one or more required fields are empty.")
	}

	return &Config{
		host:          host,
		port:          port,
		username:      username,
		password:      password,
		dbname:        dbname,
		kafka_port:    kafka_port,
		kafka_ui_port: kafka_ui_port,
		WorkerCount:   2,
		BatchSize:     5,
		Timeout:       2 * time.Second,
	}
}

func (c *Config) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.host, c.port, c.username, c.password, c.dbname)
}

// Host returns host
func (c *Config) Host() string {
	return c.host
}

// Port returns port
func (c *Config) Port() string {
	return c.port
}

// Username returns username
func (c *Config) Username() string {
	return c.username
}

// Password returns password
func (c *Config) Password() string {
	return c.password
}

// DBName returnds db name
func (c *Config) DBName() string {
	return c.dbname
}

// KafkaPort returns kafka port
func (c *Config) KafkaPort() string {
	return c.kafka_port
}

// KafkaUIPort returns kafka ui port
func (c *Config) KafkaUIPort() string {
	return c.kafka_ui_port
}

// GetRootDir returns root directory of a project
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

// ReadFirstFileWord reads first word from a file
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
