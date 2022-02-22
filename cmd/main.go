package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/xsadia/secred/api"
)

func getEnv(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv(key)
}

func main() {
	s := api.Server{}
	s.InitializeRoutes()
	s.InitializeDB(
		getEnv("APP_DB_HOST"),
		getEnv("APP_DB_USERNAME"),
		getEnv("APP_DB_PASSWORD"),
		getEnv("APP_DB_NAME"),
	)
	s.Run(":1337")
}
