package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/xsadia/secred/cmd/app"
)

func getEnv(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv(key)
}

func main() {
	a := app.App{}
	a.Initialize(
		getEnv("APP_DB_HOST"),
		getEnv("APP_DB_USERNAME"),
		getEnv("APP_DB_PASSWORD"),
		getEnv("APP_DB_NAME"),
	)
	a.Run(":1337")
}
