package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(host, user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()

	a.initializeRoutes()
}

func (a *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		response := make(map[string]bool)

		response["ok"] = true

		w.Header().Set("Content-type", "application/json")
		jsonResp, _ := json.Marshal(&response)

		w.Write(jsonResp)

	}).Methods("GET")
}
