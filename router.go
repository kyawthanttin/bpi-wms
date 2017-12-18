package main

import (
	"log"

	"github.com/kyawthanttin/bpi-wms/authentication"
	"github.com/kyawthanttin/bpi-wms/category"
	"github.com/kyawthanttin/bpi-wms/config"
	"github.com/kyawthanttin/bpi-wms/country"
	"github.com/kyawthanttin/bpi-wms/customer"
	"github.com/kyawthanttin/bpi-wms/dbutil"
	"github.com/kyawthanttin/bpi-wms/item"
	"github.com/kyawthanttin/bpi-wms/supplier"
	"github.com/kyawthanttin/bpi-wms/unitofmeasurement"
	"github.com/kyawthanttin/bpi-wms/user"
	"github.com/kyawthanttin/bpi-wms/validation"

	"github.com/gorilla/mux"
)

// NewRouter configures a new router to the API
func NewRouter() *mux.Router {
	db, err := dbutil.NewDB(config.DatasourceName)
	if err != nil {
		log.Panic(err)
	}

	validate := validation.NewValidator()

	env := &config.Env{DB: db, Validate: validate}

	r := mux.NewRouter().StrictSlash(true)

	r.Handle("/signin", authentication.SignIn(env)).Methods("POST")

	r.Handle("/users", authentication.Authenticate(env, user.UserList(env))).Methods("GET")
	r.Handle("/users/create", authentication.Authenticate(env, user.UserCreate(env))).Methods("POST")
	r.Handle("/users/{id:[0-9]+}", authentication.Authenticate(env, user.UserShow(env))).Methods("GET")
	r.Handle("/users/{id:[0-9]+}", authentication.Authenticate(env, user.UserUpdate(env))).Methods("PUT")
	r.Handle("/users/{id:[0-9]+}/changepassword", authentication.Authenticate(env, user.PasswordChange(env))).Methods("PUT")
	r.Handle("/users/{id:[0-9]+}", authentication.Authenticate(env, user.UserDelete(env))).Methods("DELETE")

	r.Handle("/categories", authentication.Authenticate(env, category.CategoryList(env))).Methods("GET")
	r.Handle("/categories/create", authentication.Authenticate(env, category.CategoryCreate(env))).Methods("POST")
	r.Handle("/categories/{id:[0-9]+}", authentication.Authenticate(env, category.CategoryShow(env))).Methods("GET")
	r.Handle("/categories/{id:[0-9]+}", authentication.Authenticate(env, category.CategoryUpdate(env))).Methods("PUT")
	r.Handle("/categories/{id:[0-9]+}", authentication.Authenticate(env, category.CategoryDelete(env))).Methods("DELETE")

	r.Handle("/countries", country.CountryList(env)).Methods("GET")
	r.Handle("/countries/{id:[0-9]+}", country.CountryShow(env)).Methods("GET")
	r.Handle("/countries/create", authentication.Authenticate(env, country.CountryCreate(env))).Methods("POST")
	r.Handle("/countries/{id:[0-9]+}", authentication.Authenticate(env, country.CountryUpdate(env))).Methods("PUT")
	r.Handle("/countries/{id:[0-9]+}", authentication.Authenticate(env, country.CountryDelete(env))).Methods("DELETE")

	r.Handle("/customers", authentication.Authenticate(env, customer.CustomerList(env))).Methods("GET")
	r.Handle("/customers/create", authentication.Authenticate(env, customer.CustomerCreate(env))).Methods("POST")
	r.Handle("/customers/{id:[0-9]+}", authentication.Authenticate(env, customer.CustomerShow(env))).Methods("GET")
	r.Handle("/customers/{id:[0-9]+}", authentication.Authenticate(env, customer.CustomerUpdate(env))).Methods("PUT")
	r.Handle("/customers/{id:[0-9]+}", authentication.Authenticate(env, customer.CustomerDelete(env))).Methods("DELETE")

	r.Handle("/items", authentication.Authenticate(env, item.ItemList(env))).Methods("GET")
	r.Handle("/items/create", authentication.Authenticate(env, item.ItemCreate(env))).Methods("POST")
	r.Handle("/items/{id:[0-9]+}", authentication.Authenticate(env, item.ItemShow(env))).Methods("GET")
	r.Handle("/items/{id:[0-9]+}", authentication.Authenticate(env, item.ItemUpdate(env))).Methods("PUT")
	r.Handle("/items/{id:[0-9]+}", authentication.Authenticate(env, item.ItemDelete(env))).Methods("DELETE")

	r.Handle("/suppliers", authentication.Authenticate(env, supplier.SupplierList(env))).Methods("GET")
	r.Handle("/suppliers/create", authentication.Authenticate(env, supplier.SupplierCreate(env))).Methods("POST")
	r.Handle("/suppliers/{id:[0-9]+}", authentication.Authenticate(env, supplier.SupplierShow(env))).Methods("GET")
	r.Handle("/suppliers/{id:[0-9]+}", authentication.Authenticate(env, supplier.SupplierUpdate(env))).Methods("PUT")
	r.Handle("/suppliers/{id:[0-9]+}", authentication.Authenticate(env, supplier.SupplierDelete(env))).Methods("DELETE")

	r.Handle("/unitofmeasurements", authentication.Authenticate(env, unitofmeasurement.UnitOfMeasurementList(env))).Methods("GET")
	r.Handle("/unitofmeasurements/create", authentication.Authenticate(env, unitofmeasurement.UnitOfMeasurementCreate(env))).Methods("POST")
	r.Handle("/unitofmeasurements/{id:[0-9]+}", authentication.Authenticate(env, unitofmeasurement.UnitOfMeasurementShow(env))).Methods("GET")
	r.Handle("/unitofmeasurements/{id:[0-9]+}", authentication.Authenticate(env, unitofmeasurement.UnitOfMeasurementUpdate(env))).Methods("PUT")
	r.Handle("/unitofmeasurements/{id:[0-9]+}", authentication.Authenticate(env, unitofmeasurement.UnitOfMeasurementDelete(env))).Methods("DELETE")

	// r.Handle("/stocks", authentication.Authenticate(stock.StockInfoList(env))).Methods("GET")
	// r.Handle("/stocks/{id:[0-9]+}", authentication.Authenticate(stock.StockInfoShow(env))).Methods("GET")

	// r.Handle("/transactionhistories", authentication.Authenticate(transactionhistory.TransactionHistoryList(env))).Methods("GET")
	// r.Handle("/transactionhistories/{id:[0-9]+}", authentication.Authenticate(transactionhistory.TransactionHistoryShow(env))).Methods("GET")

	return r
}
