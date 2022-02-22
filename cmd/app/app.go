package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Run(address string) {
	a.Router = mux.NewRouter()
	a.initializeRoutes()
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
