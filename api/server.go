package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/xsadia/secred/storage"
)

type Server struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *Server) InitializeDB(host, user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)

	var err error
	a.DB, err = storage.GetDB(connectionString)

	if err != nil {
		log.Fatal(err)
	}
}

func (a *Server) Run(address string) {
	log.Fatal(http.ListenAndServe(address, a.Router))
}

func (a *Server) InitializeRoutes() {

	a.Router = mux.NewRouter()

	a.Router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		response := make(map[string]bool)

		response["ok"] = true

		w.Header().Set("Content-type", "application/json")
		jsonResp, _ := json.Marshal(&response)

		w.Write(jsonResp)

	}).Methods("GET")
}
