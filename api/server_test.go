package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/xsadia/secred/config"
	"github.com/xsadia/secred/repository"
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

	var userAuthStr = []byte(`{
		"email":"testuser@example.com",
		"password":"123123"	
		}`)

	r, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(userCreationStr))
	r.Header.Set("Content-Type", "application/json")

	executeRequest(r)

	t.Run("should return error and http code 403 if password and e-mail match but account isn't activated",
		func(t *testing.T) {

			r, _ := http.NewRequest("POST", "/auth", bytes.NewBuffer(userAuthStr))
			r.Header.Set("Content-Type", "application/json")

			response := executeRequest(r)

			checkResponseCode(t, http.StatusForbidden, response.Code)

			var m map[string]string

			json.Unmarshal(response.Body.Bytes(), &m)

			if m["error"] != "Account not yet activated." {
				t.Errorf("expected error '%v', got '%v'", "Account not yet activated.", m["error"])
			}
		})

	t.Run("should return token if e-mail and password match and account is active", func(t *testing.T) {

		u := repository.User{Email: "testuser@example.com"}

		u.GetUserByEmail(s.DB)

		var emptyBody = []byte(`
			{}
		`)

		urlString := fmt.Sprintf("/user/confirm/%s", u.Id)

		r, _ := http.NewRequest("GET", urlString, bytes.NewBuffer(emptyBody))

		executeRequest(r)

		r, _ = http.NewRequest("POST", "/auth", bytes.NewBuffer(userAuthStr))
		r.Header.Set("Content-Type", "application/json")

		response := executeRequest(r)

		checkResponseCode(t, 200, response.Code)

		var m map[string]interface{}

		json.Unmarshal(response.Body.Bytes(), &m)

		_, ok := m["token"]

		if !ok {
			t.Error("Expected token got nothing")
		}

		if m["user"].(map[string]interface{})["id"] == "" {
			t.Error("expected id to be set")
		}

		if !reflect.DeepEqual(m["user"].(map[string]interface{})["email"], "testuser@example.com") {
			t.Errorf("expected e-mail to be '%v', got '%v'", "testuser@example.com", m["email"])
		}

		if !reflect.DeepEqual(m["user"].(map[string]interface{})["username"], "testUser") {
			t.Errorf("expected username to be '%v', got '%v'", "testUser", m["username"])
		}
	})

	t.Run("should return error and http code 401 if password or e-mail don't match", func(t *testing.T) {

		var userWrongPassword = []byte(`{
		"email":"testuser@example.com",
		"password":"123122"	
		}`)

		r, _ := http.NewRequest("POST", "/auth", bytes.NewBuffer(userWrongPassword))
		r.Header.Set("Content-Type", "application/json")

		response := executeRequest(r)

		checkResponseCode(t, 401, response.Code)

		var m map[string]string
		json.Unmarshal(response.Body.Bytes(), &m)

		if !reflect.DeepEqual(m["error"], wrongEmailPasswordCombinationError) {
			t.Errorf("expected '%v', got '%v'", wrongEmailPasswordCombinationError, m["error"])
		}

	})

}

func TestUserActivation(t *testing.T) {
	clearTables()

	var userCreationStr = []byte(`{
		"email":"testuser@example.com",
		"username":"testUser",
		"password":"123123"
		}`)

	r, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(userCreationStr))
	r.Header.Set("Content-Type", "application/json")

	executeRequest(r)

	u := repository.User{Email: "testuser@example.com"}
	u.GetUserByEmail(s.DB)

	t.Run("should return 308 if account activation is successful", func(t *testing.T) {

		r, _ := http.NewRequest("GET", "/user/confirm/"+u.Id, bytes.NewBuffer([]byte(`{}`)))

		response := executeRequest(r)

		checkResponseCode(t, 308, response.Code)
	})

	t.Run("should return error and http code 409 if user is already active", func(t *testing.T) {

		r, _ := http.NewRequest("GET", "/user/confirm/"+u.Id, bytes.NewBuffer([]byte(`{}`)))

		response := executeRequest(r)

		checkResponseCode(t, 409, response.Code)

		var m map[string]string

		json.Unmarshal(response.Body.Bytes(), &m)

		if m["error"] != "Account already activated" {
			t.Errorf("expected '%v', got '%v'", "Account already activated", m["error"])
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
