package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(path, key string) string {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv(key)
}
