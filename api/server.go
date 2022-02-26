package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/xsadia/secred/repository"
	"github.com/xsadia/secred/storage"
)

type Server struct {
	Router *mux.Router
	DB     *sql.DB
}

const emailAlreadyInUserError = "e-mail already in use"

func (s *Server) InitializeDB(host, user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)

	var err error
	s.DB, err = storage.NewDB(connectionString)

	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) Run(address string) {
	log.Fatal(http.ListenAndServe(address, s.Router))
}

func (s *Server) InitializeRoutes() {

	s.Router = mux.NewRouter()
	s.Router.HandleFunc("/user", s.createUser).Methods("POST")
	s.Router.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		response := make(map[string]bool)

		response["ok"] = true

		w.Header().Set("Content-type", "application/json")
		jsonResp, _ := json.Marshal(&response)

		w.Write(jsonResp)

	}).Methods("GET")
}

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var u repository.User
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if err := u.Create(s.DB); err != nil {
		if err.Error() == emailAlreadyInUserError {
			respondWithError(w, http.StatusConflict, emailAlreadyInUserError)
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusNoContent, map[string]string{})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
