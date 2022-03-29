package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
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
	invalidClaimError                  = "Invalid token claim"
	malformedJWTError                  = "Malformed JWT"
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
	ao := handlers.AllowedOrigins([]string{"*"})
	am := handlers.AllowedMethods([]string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH"})
	ah := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Authorization"})
	log.Fatal(http.ListenAndServe(address, handlers.CORS(ao, am, ah)(s.Router)))
}

func (s *Server) InitializeRoutes() {

	s.Router = mux.NewRouter()
	s.Router.HandleFunc("/user", s.CreateUserHandler).Methods("POST")
	s.Router.HandleFunc("/user/me", s.MeHandler).Methods("GET")
	s.Router.HandleFunc(
		"/user/confirm/{id:"+uuidRegexp+"}",
		s.ActivateUserHandler,
	).Methods("GET")
	s.Router.HandleFunc("/auth", s.AuthUserHandler).Methods("POST")
	s.Router.HandleFunc("/warehouse", s.GetWareHouseItemsHandler).Methods("GET")
	s.Router.HandleFunc("/warehouse", s.CreateWarehouseItemHandler).Methods("POST")
	s.Router.HandleFunc("/warehouse/{id:"+uuidRegexp+"}", s.GetWareHouseItemHandler).Methods("GET")
	s.Router.HandleFunc("/warehouse/{id:"+uuidRegexp+"}", s.UpdateWarehouseItemHandler).Methods("PATCH")
	s.Router.HandleFunc("/warehouse/{id:"+uuidRegexp+"}", s.DeleteWarehouseItemHandler).Methods("DELETE")
}
