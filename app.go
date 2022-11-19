package main

import (
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

func (a *App) Initializer(user, password, dbname string) {
  connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

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

func respondWithError(w http.ResponseWriter, code int, message string)  {
  respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{})  {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)

  json.NewEncoder(w).Encode(payload)
}
