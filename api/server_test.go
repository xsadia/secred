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
		config.GetEnvTest("APP_DB_HOST"),
		config.GetEnvTest("APP_DB_USERNAME"),
		config.GetEnvTest("APP_DB_PASSWORD"),
		config.GetEnvTest("APP_DB_NAME"),
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
		response := executeRequest(r)

		checkResponseCode(t, 409, response.Code)
		var m map[string]string
		json.Unmarshal(response.Body.Bytes(), &m)

		if m["error"] != emailAlreadyInUserError {
			t.Errorf("expected '%v', got '%v'", emailAlreadyInUserError, m["error"])
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
