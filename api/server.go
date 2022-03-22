package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/xsadia/secred/internal"
	"github.com/xsadia/secred/repository"
	"github.com/xsadia/secred/storage"
	"golang.org/x/crypto/bcrypt"
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
	s.Router.HandleFunc("/user", s.createUser).Methods("POST")
	s.Router.HandleFunc(
		"/user/confirm/{id:"+uuidRegexp+"}",
		s.activateUser,
	).Methods("GET")
	s.Router.HandleFunc("/auth", s.authUser).Methods("POST")
	s.Router.HandleFunc("/warehouse", s.getWareHouseItems).Methods("GET")
	s.Router.HandleFunc("/warehouse", s.createWarehouseItem).Methods("POST")
	s.Router.HandleFunc("/warehouse/{id:"+uuidRegexp+"}", s.getWareHouseItem).Methods("GET")
}

func (s *Server) activateUser(w http.ResponseWriter, r *http.Request) {
	var u repository.User
	var err error

	vars := mux.Vars(r)

	u.Id = vars["id"]
	err = u.GetUserById(s.DB)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if u.Activated {
		respondWithError(w, http.StatusConflict, "Account already activated")
		return
	}

	err = u.Activate(s.DB)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, internalServerError)
		return
	}

	http.Redirect(w, r, "http://google.com", http.StatusPermanentRedirect)
}

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var u repository.User
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, invalidRequestPayloadError)
		return
	}

	defer r.Body.Close()

	if err := u.GetUserByEmail(s.DB); err == nil {
		respondWithError(w, http.StatusConflict, emailAlreadyInUserError)
		return
	}

	u.Password = internal.HashPassword([]byte(u.Password), 8)

	if err := u.Create(s.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go internal.SendConfirmationEmail([]string{u.Email}, u.Id)

	respondWithJSON(w, http.StatusNoContent, map[string]string{})
}

func (s *Server) authUser(w http.ResponseWriter, r *http.Request) {
	var u repository.User
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, invalidRequestPayloadError)
		return
	}

	defer r.Body.Close()

	unHashedPassword := u.Password

	if err := u.GetUserByEmail(s.DB); err != nil {
		respondWithError(w, http.StatusUnauthorized, wrongEmailPasswordCombinationError)
		return
	}

	if err :=
		bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(unHashedPassword)); err != nil {
		respondWithError(w, http.StatusUnauthorized, wrongEmailPasswordCombinationError)
		return
	}

	if !u.Activated {
		respondWithError(w, http.StatusForbidden, "Account not yet activated.")
		return
	}

	token, _ := internal.CreateToken(u.Id, 9999)

	user := repository.User{
		Id:           u.Id,
		Email:        u.Email,
		Username:     u.Username,
		RefreshToken: u.RefreshToken,
		Activated:    u.Activated,
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{"token": token, "user": user})
}

func (s *Server) getWareHouseItems(w http.ResponseWriter, r *http.Request) {
	ah := r.Header.Get("Authorization")

	token, err := validateAuthHeader(ah)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	u, err := internal.ExtractUser(token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err = u.GetUserById(s.DB); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token claim")
		return
	}

	var count, start int
	startQuery := r.URL.Query().Get("start")
	countQuery := r.URL.Query().Get("count")

	if startQuery == "" {
		start = 0
	} else {
		start, _ = strconv.Atoi(startQuery)
	}

	if countQuery == "" {
		count = 10
	} else {
		count, _ = strconv.Atoi(countQuery)
	}

	items, err := repository.GetWarehouseItems(s.DB, start, count)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, internalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, items)
}

func (s *Server) getWareHouseItem(w http.ResponseWriter, r *http.Request) {
	ah := r.Header.Get("Authorization")

	token, err := validateAuthHeader(ah)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	u, err := internal.ExtractUser(token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err = u.GetUserById(s.DB); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token claim")
		return
	}

	var wi repository.WarehouseItem

	vars := mux.Vars(r)

	wi.Id = vars["id"]

	if err = wi.GetWarehouseItemById(s.DB); err != nil {
		respondWithError(w, http.StatusNotFound, "Item not found")
		return
	}

	respondWithJSON(w, http.StatusOK, wi)
}

func (s *Server) createWarehouseItem(w http.ResponseWriter, r *http.Request) {
	ah := r.Header.Get("Authorization")

	token, err := validateAuthHeader(ah)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	u, err := internal.ExtractUser(token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err = u.GetUserById(s.DB); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token claim")
		return
	}

	var wi repository.WarehouseItem

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&wi); err != nil {
		respondWithError(w, http.StatusBadRequest, invalidRequestPayloadError)
		return
	}

	defer r.Body.Close()

	if err := wi.CreateWarehouseItem(s.DB); err != nil {
		respondWithError(w, http.StatusConflict, "Item already registered")
		return
	}

	respondWithJSON(w, http.StatusCreated, wi)
}

func validateAuthHeader(header string) (string, error) {
	if len(header) < 1 {
		return "", errors.New("authorization header missing")
	}

	splitHeader := strings.Split(header, " ")

	if len(splitHeader) == 1 {
		return "", errors.New("invalid authorization header")
	}

	token := splitHeader[1]

	return token, nil
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
