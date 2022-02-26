package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvTest(key string) string {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv(key)
}

func GetEnv(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv(key)
}
