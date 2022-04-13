package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/xsadia/secred/internal"
	"github.com/xsadia/secred/repository"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) ActivateUserHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) MeHandler(w http.ResponseWriter, r *http.Request) {
	ah := r.Header.Get("Authorization")

	token, err := validateAuthHeader(ah)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	claims, err := internal.VerifyToken(token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var u repository.User

	u.Id = fmt.Sprintf("%v", claims["user_id"])

	err = u.GetUserById(s.DB)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, internalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (s *Server) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
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

	respondWithJSON(w, http.StatusCreated, nil)
}

func (s *Server) AuthUserHandler(w http.ResponseWriter, r *http.Request) {
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

	token, _ := internal.CreateToken(u.Id, 60*24*7)

	user := repository.User{
		Id:           u.Id,
		Email:        u.Email,
		Username:     u.Username,
		RefreshToken: u.RefreshToken,
		Activated:    u.Activated,
	}

	respondWithJSON(w, http.StatusOK, map[string]any{"token": token, "user": user})
}

func (s *Server) GetWareHouseItemsHandler(w http.ResponseWriter, r *http.Request) {
	ah := r.Header.Get("Authorization")

	token, err := validateAuthHeader(ah)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = internal.VerifyToken(token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
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

func (s *Server) GetWareHouseItemHandler(w http.ResponseWriter, r *http.Request) {
	ah := r.Header.Get("Authorization")

	token, err := validateAuthHeader(ah)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = internal.VerifyToken(token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
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

func (s *Server) CreateWarehouseItemHandler(w http.ResponseWriter, r *http.Request) {
	ah := r.Header.Get("Authorization")

	token, err := validateAuthHeader(ah)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = internal.VerifyToken(token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
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

func (s *Server) UpdateWarehouseItemHandler(w http.ResponseWriter, r *http.Request) {
	ah := r.Header.Get("Authorization")

	token, err := validateAuthHeader(ah)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = internal.VerifyToken(token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	vars := mux.Vars(r)

	var wi repository.WarehouseItem

	decoder := json.NewDecoder(r.Body)

	if err = decoder.Decode(&wi); err != nil {
		respondWithError(w, http.StatusBadRequest, invalidRequestPayloadError)
		return
	}

	defer r.Body.Close()

	wi.Id = vars["id"]

	if err = wi.GetWarehouseItemById(s.DB); err != nil {
		respondWithError(w, http.StatusNotFound, "Item not found")
		return
	}

	if err = wi.UpdateWarehouseItem(s.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, internalServerError)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (s *Server) DeleteWarehouseItemHandler(w http.ResponseWriter, r *http.Request) {
	ah := r.Header.Get("Authorization")

	token, err := validateAuthHeader(ah)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = internal.VerifyToken(token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	vars := mux.Vars(r)

	var wi repository.WarehouseItem

	wi.Id = vars["id"]

	if err = wi.GetWarehouseItemById(s.DB); err != nil {
		respondWithError(w, http.StatusNotFound, "Item not found")
		return
	}

	if err = wi.DeleteWarehouseItem(s.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, internalServerError)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

func (s *Server) UploadCSVWarehouse(w http.ResponseWriter, r *http.Request) {
	ah := r.Header.Get("Authorization")

	token, err := validateAuthHeader(ah)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = internal.VerifyToken(token)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("csvFile")

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err = internal.CheckMIMEContentType(header.Header); err != nil {
		respondWithError(w, http.StatusUnsupportedMediaType, err.Error())
		return
	}

	defer file.Close()

	tempFile, err := ioutil.TempFile("./temp", "upload-*.csv")

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	tempFile.Write(fileBytes)

	wil, err := internal.ParseCSV(tempFile.Name())

	if err != nil {
		respondWithError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	for _, wi := range wil {
		go func(curr repository.WarehouseItem) {
			curr.UpSertWarehouseItem(s.DB)
		}(wi)
	}

	respondWithJSON(w, http.StatusNoContent, nil)
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

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
