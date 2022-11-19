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
  Contacted bool
}

func (c *customer) get(db *sql.DB) error {
  return db.QueryRow("SELECT id, name, role, email, phone, contacted FROM customers WHERE id = $1", c.ID).Scan(&c.ID, &c.Name, &c.Role, &c.Email, &c.Phone, &c.Contacted)
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

    if err := rows.Scan(&c.ID, &c.Name, &c.Role, &c.Email, &c.Phone, &c.Contacted); err != nil {
      return nil, err
    }

    customers = append(customers, c)
  }

  return customers, err
}

func (c *customer) create(db *sql.DB) error {
  err := db.QueryRow("INSERT INTO customers(name, role, email, phone, contacted) VALUES ($1, $2, $3, $4, $5) RETURNING id", c.Name, c.Role, c.Email, c.Phone, c.Contacted).Scan(&c.ID)

  if err != nil {
    return err
  }

  return nil
}

func (c *customer) update(db *sql.DB) error {
  _, err := db.Exec("UPDATE customers SET name=$1, role=$2, email=$3, phone=$4, contacted=$5 WHERE id = $6", c.Name, c.Role, c.Email, c.Phone, c.Contacted, c.ID)

  if err != nil {
    return err
  }

  return nil
}
