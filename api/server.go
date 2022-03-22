package api

import (
	"database/sql"
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

const (
	emailAlreadyInUserError            = "e-mail already in use"
	wrongEmailPasswordCombinationError = "Wrong e-mail/password combination"
	internalServerError                = "Internal server error"
	invalidRequestPayloadError         = "Invalid request payload"
	uuidRegexp                         = "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$"
)

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
	s.Router.HandleFunc("/user", s.CreateUser).Methods("POST")
	s.Router.HandleFunc(
		"/user/confirm/{id:"+uuidRegexp+"}",
		s.ActivateUser,
	).Methods("GET")
	s.Router.HandleFunc("/auth", s.AuthUser).Methods("POST")
	s.Router.HandleFunc("/warehouse", s.GetWareHouseItems).Methods("GET")
	s.Router.HandleFunc("/warehouse", s.CreateWarehouseItem).Methods("POST")
	s.Router.HandleFunc("/warehouse/{id:"+uuidRegexp+"}", s.GetWareHouseItem).Methods("GET")
}
