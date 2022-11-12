package main

import (
	"database/sql"
	"errors"
)

type customer struct {
  ID int64 `json:"id"`
  Name string
  Role string
  Email string
  Phone string
}

func getCustomers(db sql.DB, start, count int) ([]customer, error)  {
  return nil, errors.New("Not implemented")
}
