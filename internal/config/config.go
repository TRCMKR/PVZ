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

var (
	errNoConfigFile = errors.New("no config file found")
)

// InitEnv inits env file from path
func InitEnv(envFile string) error {
	err := godotenv.Overload(envFile)
	if err != nil {
		return errNoConfigFile
	}

	return nil
}

// Config is a structure that contains all configuration parameters
type Config struct {
	host        string
	port        string
	username    string
	password    string
	dbname      string
	kafkaHost   string
	kafkaPort   string
	kafkaUIPort string
	appEnv      string
	WorkerCount int
	BatchSize   int
	Timeout     time.Duration
}

// NewConfig creates instance of Config
func NewConfig() Config {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	kafkaHost := os.Getenv("KAFKA_HOST")
	kafkaPort := os.Getenv("KAFKA_PORT")
	kafkaUIPort := os.Getenv("KAFKA_UI_PORT")
	appEnv := os.Getenv("APP_ENV")

	if host == "" || port == "" || username == "" || password == "" || dbname == "" ||
		kafkaHost == "" || kafkaPort == "" || kafkaUIPort == "" || appEnv == "" {
		log.Fatal("Database configuration missing: one or more required fields are empty.")
	}

	return Config{
		host:        host,
		port:        port,
		username:    username,
		password:    password,
		dbname:      dbname,
		kafkaHost:   kafkaHost,
		kafkaPort:   kafkaPort,
		kafkaUIPort: kafkaUIPort,
		appEnv:      appEnv,
		WorkerCount: 2,
		BatchSize:   5,
		Timeout:     2 * time.Second,
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

// DBName returns db name
func (c *Config) DBName() string {
	return c.dbname
}

// KafkaHost returns kafka port
func (c *Config) KafkaHost() string {
	return c.kafkaHost
}

// KafkaPort returns kafka port
func (c *Config) KafkaPort() string {
	return c.kafkaPort
}

// KafkaUIPort returns kafka ui port
func (c *Config) KafkaUIPort() string {
	return c.kafkaUIPort
}

// AppEnv returns env in which app is run
func (c *Config) AppEnv() string {
	return c.appEnv
}

// IsEmpty checks if config is empty
func (c *Config) IsEmpty() bool {
	return c.host == ""
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
