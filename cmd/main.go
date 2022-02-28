package main

import (
	"github.com/xsadia/secred/api"
	"github.com/xsadia/secred/config"
)

func main() {
	s := api.Server{}
	s.InitializeRoutes()
	s.InitializeDB(
		config.GetEnv(".env", "APP_DB_HOST"),
		config.GetEnv(".env", "APP_DB_USERNAME"),
		config.GetEnv(".env", "APP_DB_PASSWORD"),
		config.GetEnv(".env", "APP_DB_NAME"),
	)
	s.Run(":1337")
}
