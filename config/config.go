package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT  string
	DB_URL string
}

func GetEnvValue(name string) string {
	value := os.Getenv(name)
	
	if value == "" {
		log.Fatalf("couldn't find %v in .env file", name)
	}

	return value
}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("couldn't load .env file")
	}
}

func LoadConfig() Config {
	user := GetEnvValue("DB_USER")
	password := GetEnvValue("DB_PASSWORD")
	name := GetEnvValue("DB_NAME")
	host := GetEnvValue("DB_HOST")
	port := GetEnvValue("DB_PORT")
	serverPort := GetEnvValue("SERVER_PORT")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, name)

	return Config{PORT: serverPort, DB_URL: dbURL}
}