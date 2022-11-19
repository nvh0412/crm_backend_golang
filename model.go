package main

import (
	"database/sql"
)

type customer struct {
  ID int `json:"id"`
  Name string
  Role string
  Email string
  Phone string
}

func (c *customer) getCustomer(db *sql.DB) error {
  return db.QueryRow("SELECT id, name, role, email, phone FROM customers WHERE id = $1", c.ID).Scan(&c.ID, &c.Name, &c.Role, &c.Email, &c.Phone)
}

func getCustomers(db *sql.DB, start, count int) ([]customer, error)  {
  rows, err := db.Query(
    "SELECT id, name, role, email, phone FROM customers LIMIT $1 OFFSET $2",
    count, start)

  if err != nil {
    return nil, err
  }

  customers := []customer{}

  for rows.Next() {
    var c customer

    if err := rows.Scan(&c.ID, &c.Name, &c.Role, &c.Email, &c.Phone); err != nil {
      return nil, err
    }

    customers = append(customers, c)
  }

  return customers, err
}
