package main

import (
	"github.com/xsadia/secred/api"
	"github.com/xsadia/secred/config"
)

func main() {
	s := api.Server{}
	s.InitializeRoutes()
	s.InitializeDB(
		config.GetEnv("APP_DB_HOST"),
		config.GetEnv("APP_DB_USERNAME"),
		config.GetEnv("APP_DB_PASSWORD"),
		config.GetEnv("APP_DB_NAME"),
	)
	s.Run(":1337")
}
