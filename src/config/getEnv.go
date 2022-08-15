package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func GetEnvVariable(key string) string {
	cwd, e := os.Getwd()
	if e != nil {
		log.Fatalf("Permission denied for get cwd command")
	}

	env := strings.ToLower(os.Getenv("ENV"))
	if env == "production" {
		return os.Getenv(key)
	}

	// load .env file
	err := godotenv.Load(cwd + "/.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
