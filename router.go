package main

import (
	"log"

	"github.com/kyawthanttin/bpi-wms/authentication"
	"github.com/kyawthanttin/bpi-wms/category"
	"github.com/kyawthanttin/bpi-wms/config"
	"github.com/kyawthanttin/bpi-wms/country"
	"github.com/kyawthanttin/bpi-wms/customer"
	"github.com/kyawthanttin/bpi-wms/dbutil"
	"github.com/kyawthanttin/bpi-wms/supplier"
	"github.com/kyawthanttin/bpi-wms/unitofmeasurement"
	"github.com/kyawthanttin/bpi-wms/user"

	"github.com/gorilla/mux"
)

// NewRouter configures a new router to the API
func NewRouter() *mux.Router {
	db, err := dbutil.NewDB(config.DatasourceName)
	if err != nil {
		log.Panic(err)
	}

	env := &config.Env{DB: db}

	r := mux.NewRouter().StrictSlash(true)

	adminRoleName := "ADMIN"

	r.Handle("/signin", authentication.SignIn(env)).Methods("POST")

	r.Handle("/users", authentication.AuthenticateWithRoles(env, user.UserList(env), []string{adminRoleName})).Methods("GET")
	r.Handle("/users/create", authentication.AuthenticateWithRoles(env, user.UserCreate(env), []string{adminRoleName})).Methods("POST")
	r.Handle("/users/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, user.UserShow(env), nil)).Methods("GET")
	r.Handle("/users/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, user.UserUpdate(env), nil)).Methods("PUT")
	r.Handle("/users/{id:[0-9]+}/changepassword", authentication.AuthenticateWithRoles(env, user.PasswordChange(env), nil)).Methods("PUT")
	r.Handle("/users/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, user.UserDelete(env), []string{adminRoleName})).Methods("DELETE")

	r.Handle("/categories", authentication.AuthenticateWithRoles(env, category.CategoryList(env), []string{adminRoleName})).Methods("GET")
	r.Handle("/categories/create", authentication.AuthenticateWithRoles(env, category.CategoryCreate(env), []string{adminRoleName})).Methods("POST")
	r.Handle("/categories/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, category.CategoryShow(env), nil)).Methods("GET")
	r.Handle("/categories/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, category.CategoryUpdate(env), []string{adminRoleName})).Methods("PUT")
	r.Handle("/categories/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, category.CategoryDelete(env), []string{adminRoleName})).Methods("DELETE")

	r.Handle("/countries", country.CountryList(env)).Methods("GET")
	r.Handle("/countries/{id:[0-9]+}", country.CountryShow(env)).Methods("GET")
	r.Handle("/countries/create", authentication.AuthenticateWithRoles(env, country.CountryCreate(env), []string{adminRoleName})).Methods("POST")
	r.Handle("/countries/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, country.CountryUpdate(env), []string{adminRoleName})).Methods("PUT")
	r.Handle("/countries/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, country.CountryDelete(env), []string{adminRoleName})).Methods("DELETE")

	r.Handle("/customers", authentication.AuthenticateWithRoles(env, customer.CustomerList(env), []string{adminRoleName})).Methods("GET")
	r.Handle("/customers/create", authentication.AuthenticateWithRoles(env, customer.CustomerCreate(env), []string{adminRoleName})).Methods("POST")
	r.Handle("/customers/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, customer.CustomerShow(env), nil)).Methods("GET")
	r.Handle("/customers/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, customer.CustomerUpdate(env), []string{adminRoleName})).Methods("PUT")
	r.Handle("/customers/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, customer.CustomerDelete(env), []string{adminRoleName})).Methods("DELETE")

	// r.Handle("/items", authentication.AuthenticateWithRoles(env, item.ItemList(env), []string{adminRoleName})).Methods("GET")
	// r.Handle("/items/create", authentication.AuthenticateWithRoles(env, item.ItemCreate(env), []string{adminRoleName})).Methods("POST")
	// r.Handle("/items/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, item.ItemShow(env), nil)).Methods("GET")
	// r.Handle("/items/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, item.ItemUpdate(env), []string{adminRoleName})).Methods("PUT")
	// r.Handle("/items/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, item.ItemDelete(env), []string{adminRoleName})).Methods("DELETE")

	r.Handle("/suppliers", authentication.AuthenticateWithRoles(env, supplier.SupplierList(env), []string{adminRoleName})).Methods("GET")
	r.Handle("/suppliers/create", authentication.AuthenticateWithRoles(env, supplier.SupplierCreate(env), []string{adminRoleName})).Methods("POST")
	r.Handle("/suppliers/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, supplier.SupplierShow(env), nil)).Methods("GET")
	r.Handle("/suppliers/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, supplier.SupplierUpdate(env), []string{adminRoleName})).Methods("PUT")
	r.Handle("/suppliers/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, supplier.SupplierDelete(env), []string{adminRoleName})).Methods("DELETE")

	r.Handle("/unitofmeasurements", authentication.AuthenticateWithRoles(env, unitofmeasurement.UnitOfMeasurementList(env), []string{adminRoleName})).Methods("GET")
	r.Handle("/unitofmeasurements/create", authentication.AuthenticateWithRoles(env, unitofmeasurement.UnitOfMeasurementCreate(env), []string{adminRoleName})).Methods("POST")
	r.Handle("/unitofmeasurements/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, unitofmeasurement.UnitOfMeasurementShow(env), nil)).Methods("GET")
	r.Handle("/unitofmeasurements/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, unitofmeasurement.UnitOfMeasurementUpdate(env), []string{adminRoleName})).Methods("PUT")
	r.Handle("/unitofmeasurements/{id:[0-9]+}", authentication.AuthenticateWithRoles(env, unitofmeasurement.UnitOfMeasurementDelete(env), []string{adminRoleName})).Methods("DELETE")

	// r.Handle("/stocks", authentication.Authenticate(stock.StockInfoList(env))).Methods("GET")
	// r.Handle("/stocks/{id:[0-9]+}", authentication.Authenticate(stock.StockInfoShow(env))).Methods("GET")

	// r.Handle("/transactionhistories", authentication.Authenticate(transactionhistory.TransactionHistoryList(env))).Methods("GET")
	// r.Handle("/transactionhistories/{id:[0-9]+}", authentication.Authenticate(transactionhistory.TransactionHistoryShow(env))).Methods("GET")

	return r
}
