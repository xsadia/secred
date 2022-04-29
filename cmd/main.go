package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/xsadia/secred/api"
)

func main() {
	godotenv.Load(".env")
	s := api.Server{}
	s.InitializeRoutes()
	s.InitializeDB(
		os.Getenv("APP_DB_HOST"),
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)
	s.Run(":1337")
}
