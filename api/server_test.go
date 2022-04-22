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

type ItemResponse struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Error    string `json:"error"`
	Quantity int    `json:"quantity"`
	Min      int    `json:"min"`
	Max      int    `json:"max"`
}

var (
	s Server

	userCreationStr = []byte(`{
	"email":"testuser@example.com", 
	"username":"testUser",
	"password":"123123"
	}`)

	userAuthStr = []byte(`{
	"email":"testuser@example.com",
	"password":"123123"	
	}`)
)

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
		checkResponseCode(t, http.StatusCreated, response.Code)
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

func TestMe(t *testing.T) {
	r, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(userCreationStr))
	r.Header.Set("Content-Type", "application/json")

	executeRequest(r)

	u := repository.User{Email: "testuser@example.com"}
	u.GetUserByEmail(s.DB)

	r, _ = http.NewRequest("GET", "/user/confirm/"+u.Id, nil)

	executeRequest(r)

	r, _ = http.NewRequest("POST", "/auth", bytes.NewBuffer(userAuthStr))

	authResponse := executeRequest(r)

	var am map[string]interface{}

	json.Unmarshal(authResponse.Body.Bytes(), &am)

	tokenString := fmt.Sprintf("bearer %v", am["token"])

	t.Run("Should return user's information if authenticated", func(t *testing.T) {
		r, _ := http.NewRequest("GET", "/user/me", nil)
		r.Header.Set("Authorization", tokenString)

		response := executeRequest(r)

		checkResponseCode(t, http.StatusOK, response.Code)

		var m map[string]interface{}

		json.Unmarshal(response.Body.Bytes(), &m)

		if m["id"] == "" {
			t.Error("expected id to be set")
		}

		if !reflect.DeepEqual(m["email"], "testuser@example.com") {
			t.Errorf("expected e-mail to be '%v', got '%v'", "testuser@example.com", m["email"])
		}

		if !reflect.DeepEqual(m["username"], "testUser") {
			t.Errorf("expected username to be '%v', got '%v'", "testUser", m["username"])
		}

	})

	t.Run("Should return error if user is not authenticated", func(t *testing.T) {
		r, _ := http.NewRequest("GET", "/user/me", nil)

		response := executeRequest(r)

		checkResponseCode(t, http.StatusUnauthorized, response.Code)

	})
}

func TestAuthUser(t *testing.T) {
	clearTables()

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

		urlString := fmt.Sprintf("/user/confirm/%s", u.Id)

		r, _ := http.NewRequest("GET", urlString, nil)

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

	r, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(userCreationStr))
	r.Header.Set("Content-Type", "application/json")

	executeRequest(r)

	u := repository.User{Email: "testuser@example.com"}
	u.GetUserByEmail(s.DB)

	t.Run("should return 308 if account activation is successful", func(t *testing.T) {

		r, _ := http.NewRequest("GET", "/user/confirm/"+u.Id, nil)

		response := executeRequest(r)

		checkResponseCode(t, 308, response.Code)
	})

	t.Run("should return error and http code 409 if user is already active", func(t *testing.T) {

		r, _ := http.NewRequest("GET", "/user/confirm/"+u.Id, nil)

		response := executeRequest(r)

		checkResponseCode(t, http.StatusConflict, response.Code)

		var m map[string]string

		json.Unmarshal(response.Body.Bytes(), &m)

		if m["error"] != "Account already activated" {
			t.Errorf("expected '%v', got '%v'", "Account already activated", m["error"])
		}
	})
}

func TestWarehouseItems(t *testing.T) {
	clearTables()

	itemCreationString := []byte(`{
		"name": "testItem",
		"quantity": 2,
		"min": 1,
		"max": 3
	}`)

	r, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(userCreationStr))
	r.Header.Set("Content-Type", "application/json")

	executeRequest(r)

	u := repository.User{Email: "testuser@example.com"}
	u.GetUserByEmail(s.DB)

	r, _ = http.NewRequest("GET", "/user/confirm/"+u.Id, nil)

	executeRequest(r)

	r, _ = http.NewRequest("POST", "/auth", bytes.NewBuffer(userAuthStr))

	authResponse := executeRequest(r)

	var am map[string]interface{}

	json.Unmarshal(authResponse.Body.Bytes(), &am)

	tokenString := fmt.Sprintf("bearer %v", am["token"])

	var rs ItemResponse

	t.Run("Should create item if user is authorized", func(t *testing.T) {

		r, _ := http.NewRequest("POST", "/warehouse", bytes.NewBuffer(itemCreationString))
		r.Header.Set("Authorization", tokenString)

		response := executeRequest(r)

		checkResponseCode(t, http.StatusCreated, response.Code)

		json.Unmarshal(response.Body.Bytes(), &rs)

		if rs.Error != "" {
			t.Errorf("Expected error to be nil got %q", rs.Error)
		}

		if rs.Name != "testItem" {
			t.Errorf("Expected name to be testItem, got %q", rs.Name)
		}

		if rs.Quantity != 2 {
			t.Errorf("Expected quantity to be 2, got %d", rs.Quantity)
		}

		if rs.Min != 1 {
			t.Errorf("Expected min to be 1, got %d", rs.Min)
		}

		if rs.Max != 3 {
			t.Errorf("Expected min to be 3, got %d", rs.Max)
		}
	})

	t.Run("Should return error if a item is already registered with the same name", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/warehouse", bytes.NewBuffer(itemCreationString))
		r.Header.Set("Authorization", tokenString)

		response := executeRequest(r)

		checkResponseCode(t, http.StatusConflict, response.Code)
	})

	t.Run("Should return items if user is authorized", func(t *testing.T) {
		r, _ := http.NewRequest("GET", "/warehouse", nil)
		r.Header.Set("Authorization", tokenString)

		response := executeRequest(r)

		checkResponseCode(t, http.StatusOK, response.Code)

		var items []ItemResponse
		json.Unmarshal(response.Body.Bytes(), &items)

		if len(items) != 1 {
			t.Errorf("Expected array to be of size 1, got %d", len(items))
		}

		if items[0].Error != "" {
			t.Errorf("Expected error to be nil got %q", items[0].Error)
		}

		if items[0].Name != "testItem" {
			t.Errorf("Expected name to be testItem, got %q", items[0].Name)
		}

		if items[0].Quantity != 2 {
			t.Errorf("Expected quantity to be 2, got %d", items[0].Quantity)
		}

		if items[0].Min != 1 {
			t.Errorf("Expected min to be 1, got %d", items[0].Min)
		}

		if items[0].Max != 3 {
			t.Errorf("Expected min to be 3, got %d", items[0].Max)
		}
	})

	t.Run("Should return item based on it's id if user is authorized", func(t *testing.T) {
		r, _ := http.NewRequest("GET", "/warehouse/"+rs.Id, nil)
		r.Header.Set("Authorization", tokenString)

		response := executeRequest(r)

		checkResponseCode(t, http.StatusOK, response.Code)

		var item ItemResponse
		json.Unmarshal(response.Body.Bytes(), &item)

		if item.Error != "" {
			t.Errorf("Expected error to be nil got %q", item.Error)
		}

		if item.Name != "testItem" {
			t.Errorf("Expected name to be testItem, got %q", item.Name)
		}

		if item.Quantity != 2 {
			t.Errorf("Expected quantity to be 2, got %d", item.Quantity)
		}

		if item.Min != 1 {
			t.Errorf("Expected min to be 1, got %d", item.Min)
		}

		if item.Max != 3 {
			t.Errorf("Expected min to be 3, got %d", item.Max)
		}
	})

	t.Run("Should update item if item exists", func(t *testing.T) {
		updateStr := []byte(`{
			"quantity": 10,
			"min": 3,
			"max": 15
		}`)

		r, _ := http.NewRequest("PATCH", "/warehouse/"+rs.Id, bytes.NewBuffer(updateStr))
		r.Header.Set("Authorization", tokenString)

		response := executeRequest(r)

		checkResponseCode(t, http.StatusNoContent, response.Code)

		r, _ = http.NewRequest("GET", "/warehouse/"+rs.Id, nil)
		r.Header.Set("Authorization", tokenString)

		response = executeRequest(r)

		checkResponseCode(t, http.StatusOK, response.Code)

		var item ItemResponse
		json.Unmarshal(response.Body.Bytes(), &item)

		if item.Error != "" {
			t.Errorf("Expected error to be nil got %q", item.Error)
		}

		if item.Name != "testItem" {
			t.Errorf("Expected name to be testItem, got %q", item.Name)
		}

		if item.Quantity != 2 {
			t.Errorf("Expected quantity to be 10, got %d", item.Quantity)
		}

		if item.Min != 1 {
			t.Errorf("Expected min to be 2, got %d", item.Min)
		}

		if item.Max != 3 {
			t.Errorf("Expected min to be 15, got %d", item.Max)
		}
	})

	t.Run("Should throw error when updating a item that doesn't exist", func(t *testing.T) {
		updateStr := []byte(`{
			"quantity": 10,
			"min": 3,
			"max": 15
		}`)

		r, _ := http.NewRequest("PATCH", "/warehouse/150e6365-c6cc-48f2-8ecf-e076a0e1e8b7", bytes.NewBuffer(updateStr))
		r.Header.Set("Authorization", tokenString)

		response := executeRequest(r)

		checkResponseCode(t, http.StatusNotFound, response.Code)
	})

	t.Run("Should delete item if item exists", func(t *testing.T) {
		r, _ := http.NewRequest("DELETE", "/warehouse/"+rs.Id, nil)
		r.Header.Set("Authorization", tokenString)

		response := executeRequest(r)

		checkResponseCode(t, http.StatusNoContent, response.Code)

		r, _ = http.NewRequest("GET", "/warehouse/"+rs.Id, nil)
		r.Header.Set("Authorization", tokenString)

		response = executeRequest(r)

		var item ItemResponse
		json.Unmarshal(response.Body.Bytes(), &item)

		if item.Error != "Item not found" {
			t.Errorf("Expected error got %v", item.Error)
		}
	})

	t.Run("Should throw error if item doesn't exist", func(t *testing.T) {
		r, _ := http.NewRequest("DELETE", "/warehouse/"+rs.Id, nil)
		r.Header.Set("Authorization", tokenString)

		response := executeRequest(r)

		checkResponseCode(t, http.StatusNotFound, response.Code)

		var item ItemResponse
		json.Unmarshal(response.Body.Bytes(), &item)

		if item.Error != "Item not found" {
			t.Errorf("Expected error got %v", item.Error)
		}
	})
}

func TestSchool(t *testing.T) {
	clearTables()
	schoolCreationString := []byte(`{
		"name": "CSC"
	}
	`)
	token := createAndAuthUser()

	t.Run("Should create school", func(t *testing.T) {
		r, _ := http.NewRequest("POST", "/school", bytes.NewBuffer(schoolCreationString))
		r.Header.Set("Authorization", token)

		response := executeRequest(r)

		checkResponseCode(t, http.StatusCreated, response.Code)

		var school repository.School
		json.Unmarshal(response.Body.Bytes(), &school)

		if school.Name != "CSC" {
			t.Errorf("Expected school name to be %q, got %q", "CSC", school.Name)
		}
	})
}

func createAndAuthUser() string {

	r, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(userCreationStr))
	r.Header.Set("Content-Type", "application/json")

	executeRequest(r)

	u := repository.User{Email: "testuser@example.com"}
	u.GetUserByEmail(s.DB)

	r, _ = http.NewRequest("GET", "/user/confirm/"+u.Id, nil)

	executeRequest(r)

	r, _ = http.NewRequest("POST", "/auth", bytes.NewBuffer(userAuthStr))

	authResponse := executeRequest(r)

	var am map[string]interface{}

	json.Unmarshal(authResponse.Body.Bytes(), &am)

	tokenString := fmt.Sprintf("bearer %v", am["token"])

	return tokenString
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

	s.DB.Exec("DELETE FROM warehouse_items")

	s.DB.Exec("DELETE FROM schools")
}
