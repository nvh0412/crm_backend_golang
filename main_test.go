package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS customers
(
  id SERIAL,
  name TEXT NOT NULL,
  role TEXT NOT NULL,
  email TEXT NOT NULL,
  phone TEXT NOT NULL,
  contacted BOOLEAN NOT NULL
)
`
func addCustomer(count int) {
  if (count < 1) {
    count = 1
  }

  for i := 0; i < count; i++ {
    a.DB.Exec("INSERT INTO customers(name, role, email, phone, contacted) VALUES ($1, $2, $3, $4, $5)", "Customer" + strconv.Itoa(i), "Role", "customer" + strconv.Itoa(i) + "@gmail.com", "012345667", true)
  }
}

func tearUp() {
  if _, err := a.DB.Exec(tableCreationQuery); err != nil {
    log.Fatal(err)
  }
}

func tearDown() {
  a.DB.Exec("DELETE FROM customers")
  a.DB.Exec("ALTER SEQUENCE customers_id_seq RESTART WITH 1")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
  rr := httptest.NewRecorder()

  a.Router.ServeHTTP(rr, req)

  return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
  if expected != actual {
    t.Errorf("Expected response code %d. Got %d\n", expected, actual)
  }
}

func TestMain(m *testing.M)  {
  a.Initializer(
    os.Getenv("APP_DB_USERNAME"),
    os.Getenv("APP_DB_PASSWORD"),
    os.Getenv("APP_DB_NAME"),
  )
  tearDown()

  tearUp()
  code := m.Run()
  tearDown()
  os.Exit(code)
}

func TestGetCustomers(t *testing.T) {
  tearDown()

  req, _ := http.NewRequest("GET", "/customers", nil)

  response := executeRequest(req)

  checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetCustomer(t *testing.T)  {
  tearDown()
  addCustomer(1)

  req, _ := http.NewRequest("GET", "/customers/1", nil)
  response := executeRequest(req)

  checkResponseCode(t, http.StatusOK, response.Code)
}

func TestCreateCustomer(t *testing.T) {
  tearDown()

  body := []byte(`{"name": "Gopher", "role": "employee", "email": "go@go.com", "phone": "123123123", "contacted": false}`)

  req, _ := http.NewRequest("POST", "/customers", bytes.NewBuffer(body))
  response := executeRequest(req)

  checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestUpdateCustomer(t *testing.T) {
  tearDown()
  addCustomer(1)

  body := []byte(`{"name": "Gopher", "role": "employee", "email": "go@go.com", "phone": "123123123", "contacted": false}`)

  req, _ := http.NewRequest("PUT", "/customers/1", bytes.NewBuffer(body))
  response := executeRequest(req)

  checkResponseCode(t, http.StatusOK, response.Code)
}

func TestDeleteCustomer(t *testing.T) {
  tearDown()

  req, _ := http.NewRequest("DELETE", "/customers/1", nil)

  response := executeRequest(req)

  checkResponseCode(t, http.StatusNotFound, response.Code)

  addCustomer(1)

  req2, _ := http.NewRequest("DELETE", "/customers/1", nil)

  response2 := executeRequest(req2)

  checkResponseCode(t, http.StatusOK, response2.Code)
}
