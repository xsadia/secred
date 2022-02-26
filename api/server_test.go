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
		config.GetEnv("APP_DB_HOST"),
		config.GetEnv("APP_DB_USERNAME"),
		config.GetEnv("APP_DB_PASSWORD"),
		config.GetEnv("APP_DB_NAME"),
	)

	code := m.Run()

	clearTables()

	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	clearTables()
	var jsonStr = []byte(`{"email":"testuser@example.com", "username":"testUser","password":"123123"}`)
	r, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	r.Header.Set("Content-Type", "application/json")

	response := executeRequest(r)
	checkResponseCode(t, 201, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["email"] != "testuser@example.com" {
		t.Errorf("expected email to be %q, got %q", "testuser@example", m["email"])
	}

	if m["username"] != "testUser" {
		t.Errorf("")
	}
}

func executeRequest(r *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.Router.ServeHTTP(rr, r)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func clearTables() {
	s.DB.Exec("DELETE FROM users")
}
