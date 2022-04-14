package api

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		internal.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if u.Activated {
		internal.RespondWithError(w, http.StatusConflict, "Account already activated")
		return
	}

	err = u.Activate(s.DB)

	if err != nil {
		internal.RespondWithError(w, http.StatusInternalServerError, internalServerError)
		return
	}

	http.Redirect(w, r, "http://google.com", http.StatusPermanentRedirect)
}

func (s *Server) MeHandler(w http.ResponseWriter, r *http.Request) {
	ah := r.Header.Get("Authorization")

	token, err := internal.ValidateAuthHeader(ah)

	if err != nil {
		internal.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	claims, err := internal.VerifyToken(token)

	if err != nil {
		internal.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var u repository.User

	u.Id = fmt.Sprintf("%v", claims["user_id"])

	err = u.GetUserById(s.DB)

	if err != nil {
		internal.RespondWithError(w, http.StatusInternalServerError, internalServerError)
		return
	}

	internal.RespondWithJSON(w, http.StatusOK, u)
}

func (s *Server) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var u repository.User
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&u); err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, invalidRequestPayloadError)
		return
	}

	defer r.Body.Close()

	if err := u.GetUserByEmail(s.DB); err == nil {
		internal.RespondWithError(w, http.StatusConflict, emailAlreadyInUserError)
		return
	}

	u.Password = internal.HashPassword([]byte(u.Password), 8)

	if err := u.Create(s.DB); err != nil {
		internal.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	go internal.SendConfirmationEmail([]string{u.Email}, u.Id)

	internal.RespondWithJSON(w, http.StatusCreated, nil)
}

func (s *Server) AuthUserHandler(w http.ResponseWriter, r *http.Request) {
	var u repository.User
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&u); err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, invalidRequestPayloadError)
		return
	}

	defer r.Body.Close()

	unHashedPassword := u.Password

	if err := u.GetUserByEmail(s.DB); err != nil {
		internal.RespondWithError(w, http.StatusUnauthorized, wrongEmailPasswordCombinationError)
		return
	}

	if err :=
		bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(unHashedPassword)); err != nil {
		internal.RespondWithError(w, http.StatusUnauthorized, wrongEmailPasswordCombinationError)
		return
	}

	if !u.Activated {
		internal.RespondWithError(w, http.StatusForbidden, "Account not yet activated.")
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

	internal.RespondWithJSON(w, http.StatusOK, map[string]any{"token": token, "user": user})
}
