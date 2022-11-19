package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS customers
(
  id SERIAL,
  name TEXT NOT NULL,
  role TEXT NOT NULL,
  email TEXT NOT NULL,
  phone TEXT NOT NULL
)
`

func tearUp() {
  if _, err := a.DB.Exec(tableCreationQuery); err != nil {
    log.Fatal(err)
  }
}

func tearDown() {
  a.DB.Exec("DELETE FROM customers")
  a.DB.Exec("ALTER SEQUENCE customers_id_seq RESTART WITH 1")
}


func TestMain(m *testing.M)  {
  a.Initializer(
    os.Getenv("APP_DB_USERNAME"),
    os.Getenv("APP_DB_PASSWORD"),
    os.Getenv("APP_DB_NAME"),
  )

  tearUp()
  code := m.Run()
  tearDown()
  os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
  tearDown()

  req, _ := http.NewRequest("GET", "/customers", nil)

  response := executeRequest(req)

  checkResponseCode(t, http.StatusOK, response.Code)
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

