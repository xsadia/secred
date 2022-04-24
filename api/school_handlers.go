package api

import (
	"encoding/json"
	"net/http"

	"github.com/xsadia/secred/internal"
	"github.com/xsadia/secred/repository"
)

func (s *Server) CreateSchoolHandler(w http.ResponseWriter, r *http.Request) {
	ah := r.Header.Get("Authorization")

	token, err := internal.ValidateAuthHeader(ah)

	if err != nil {
		internal.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	_, err = internal.VerifyToken(token)

	if err != nil {
		internal.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	var sc repository.School

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&sc); err != nil {
		internal.RespondWithError(w, http.StatusBadRequest, invalidRequestPayloadError)
		return
	}

	defer r.Body.Close()

	if err = sc.CreateSchool(s.DB); err != nil {
		internal.RespondWithError(w, http.StatusConflict, "School already registered")
		return
	}

	internal.RespondWithJSON(w, http.StatusCreated, sc)
}
