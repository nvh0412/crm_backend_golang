package main

import (
	"log"
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
