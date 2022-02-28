package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/xsadia/secred/config"
)

var s Server

func TestMain(m *testing.M) {
	s = Server{}

	s.InitializeRoutes()
	s.InitializeDB(
		config.GetEnv("../.env", "APP_DB_HOST"),
		config.GetEnv("../.env", "APP_DB_USERNAME"),
		config.GetEnv("../.env", "APP_DB_PASSWORD"),
		config.GetEnv("../.env", "APP_DB_NAME"),
	)

	code := m.Run()

	clearTables()

	os.Exit(code)
}

func TestCreateUser(t *testing.T) {

	t.Run("should create an user and return http code 204 if e-mail is available", func(t *testing.T) {
		clearTables()

		var jsonStr = []byte(`{
		"email":"testuser@example.com", 
		"username":"testUser",
		"password":"123123"
		}`)
		r, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
		r.Header.Set("Content-Type", "application/json")

		response := executeRequest(r)
		checkResponseCode(t, 204, response.Code)
	})

	t.Run("should not create user and return http code 409 if e-mail is already in use", func(t *testing.T) {
		clearTables()

		var jsonStr = []byte(`{
		"email":"testuser@example.com", 
		"username":"testUser",
		"password":"123123"
		}`)
		r, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
		r.Header.Set("Content-Type", "application/json")

		executeRequest(r)

		r, _ = http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
		r.Header.Set("Content-Type", "application/json")
		response := executeRequest(r)

		checkResponseCode(t, 409, response.Code)
		var m map[string]string
		json.Unmarshal(response.Body.Bytes(), &m)

		if m["error"] != emailAlreadyInUserError {
			t.Errorf("expected '%v', got '%v'", emailAlreadyInUserError, m["error"])
		}
	})
}

func TestAuthUser(t *testing.T) {
	clearTables()

	var userCreationStr = []byte(`{
		"email":"testuser@example.com", 
		"username":"testUser",
		"password":"123123"
		}`)

	r, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(userCreationStr))
	r.Header.Set("Content-Type", "application/json")

	executeRequest(r)

	t.Run("should return token if e-mail and password match", func(t *testing.T) {

		var userAuthStr = []byte(`{
			"email":"testuser@example.com", 
			"password":"123123"
			}`)

		r, _ := http.NewRequest("POST", "/auth", bytes.NewBuffer(userAuthStr))
		r.Header.Set("Content-Type", "application/json")

		response := executeRequest(r)

		checkResponseCode(t, 200, response.Code)

		var m map[string]string

		json.Unmarshal(response.Body.Bytes(), &m)

		_, ok := m["token"]

		if !ok {
			t.Error("Expected token got nothing")
		}
	})

	t.Run("should return error and http code 401 if password or e-mail don't match", func(t *testing.T) {

		var userAuthStr = []byte(`{
			"email":"testuser@example.com", 
			"password":"123122"
			}`)

		r, _ := http.NewRequest("POST", "/auth", bytes.NewBuffer(userAuthStr))
		r.Header.Set("Content-Type", "application/json")

		response := executeRequest(r)

		checkResponseCode(t, 401, response.Code)

		var m map[string]string
		json.Unmarshal(response.Body.Bytes(), &m)

		if m["error"] != wrongEmailPasswordCombinationError {
			t.Errorf("expected '%v', got '%v'", wrongEmailPasswordCombinationError, m["error"])
		}
	})
}

func executeRequest(r *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, r)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	t.Helper()
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func clearTables() {
	s.DB.Exec("DELETE FROM users")
}
