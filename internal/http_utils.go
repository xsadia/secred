package internal

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func ValidateAuthHeader(header string) (string, error) {
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

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
