package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
  Router *mux.Router
  DB *sql.DB
}

func (a *App) Initializer(user, password, dbname, connStr string) {
  var connectionString string

  if connStr != "" {
    connectionString = fmt.Sprintf("%s?sslmode=disable", connStr)
  } else {
    connectionString = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
  }

  var err error
  a.DB, err = sql.Open("postgres", connectionString)

  if err != nil {
    log.Fatal(err)
  }

  a.Router = mux.NewRouter()

  a.InitializeRoute()
}

func (a *App) Run(addr string)  {
  log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) InitializeRoute()  {
  a.Router.HandleFunc("/customers", a.getCustomers).Methods("GET")
  a.Router.HandleFunc("/customers/{id:[0-9]+}", a.getCustomer).Methods("GET")
  a.Router.HandleFunc("/customers", a.createCustomer).Methods("POST")
  a.Router.HandleFunc("/customers/{id:[0-9]+}", a.updateCustomer).Methods("PUT")
  a.Router.HandleFunc("/customers/{id:[0-9]+}", a.destroyCustomer).Methods("DELETE")
  a.Router.HandleFunc("/customers/bulk", a.updateCustomersInBatch).Methods("PUT")
}

func (a *App) getCustomer(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  id, err := strconv.Atoi(vars["id"])

  if err != nil {
    respondWithError(w, http.StatusBadRequest, err.Error())
  }

  customer := customer{ID: id}

  if err := customer.get(a.DB); err != nil {
    respondWithError(w, http.StatusNotFound, err.Error())
  }

  respondWithJSON(w, http.StatusOK, customer)
}

func (a *App) getCustomers(w http.ResponseWriter, r *http.Request) {
  count, _ := strconv.Atoi(r.FormValue("count"))
  start, _ := strconv.Atoi(r.FormValue("start"))

  // Default values
  if count > 10 || count < 1 {
    count = 10
  }

  if start < 0 {
    start = 0
  }

  customers, err := getCustomers(a.DB, start, count)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
  }

  respondWithJSON(w, http.StatusOK, customers)
}

func (a *App) createCustomer(w http.ResponseWriter, r *http.Request) {
  var c customer

  decoder := json.NewDecoder(r.Body)
  if err := decoder.Decode(&c); err != nil {
    respondWithError(w, http.StatusBadRequest, "Invalid params")
    return
  }

  if err := c.create(a.DB); err != nil {
    fmt.Println(err)
    respondWithError(w, http.StatusUnprocessableEntity, "Unprocessable Entity")
    return
  }

  respondWithJSON(w, http.StatusCreated, c)
}

func (a *App) updateCustomersInBatch(w http.ResponseWriter, r *http.Request) {
  var customers []customer

  decoder := json.NewDecoder(r.Body)
  if err := decoder.Decode(&customers); err != nil {
    respondWithError(w, http.StatusBadRequest, "Invalid params")
    return
  }

  ctx := context.Background()
  tx, err := a.DB.BeginTx(ctx, nil)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
  }

  defer tx.Rollback()

  for _, customer := range customers {
    if err := customer.update(a.DB); err != nil {
      tx.Rollback()
      respondWithError(w, http.StatusUnprocessableEntity, err.Error())
      return
    }
  }

  if err = tx.Commit(); err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
  }

  respondWithJSON(w, http.StatusOK, customers)
}

func (a *App) updateCustomer(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  id, err := strconv.Atoi(vars["id"])

  if err != nil {
    respondWithError(w, http.StatusBadRequest, err.Error())
  }

  c := customer{ID: id}

  if err := c.get(a.DB); err != nil {
    respondWithError(w, http.StatusNotFound, err.Error())
  }

  decoder := json.NewDecoder(r.Body)
  if err := decoder.Decode(&c); err != nil {
    respondWithError(w, http.StatusBadRequest, "Invalid params")
    return
  }

  if err := c.update(a.DB); err != nil {
    respondWithError(w, http.StatusUnprocessableEntity, "Unprocessable Entity")
    return
  }

  respondWithJSON(w, http.StatusOK, c)
}

func (a *App) destroyCustomer(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  id, err := strconv.Atoi(vars["id"])

  if err != nil {
    respondWithError(w, http.StatusBadRequest, err.Error())
  }

  c := customer{ID: id}

  if err := c.get(a.DB); err != nil {
    respondWithError(w, http.StatusNotFound, err.Error())
  }

  if err := c.destroy(a.DB); err != nil {
    respondWithError(w, http.StatusUnprocessableEntity, "Unprocessable Entity")
    return
  }

  respondWithJSON(w, http.StatusOK, c)
}

func respondWithError(w http.ResponseWriter, code int, message string)  {
  respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{})  {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)

  json.NewEncoder(w).Encode(payload)
}
