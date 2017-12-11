package main

import (
	"log"

	"github.com/kyawthanttin/bpi-wms/authentication"
	"github.com/kyawthanttin/bpi-wms/config"

	"github.com/gorilla/mux"
	"github.com/kyawthanttin/bpi-wms/category"
)

// NewRouter configures a new router to the API
func NewRouter() *mux.Router {
	db, err := config.NewDB("user=wms password=wms dbname=bpi_wms sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	env := &config.Env{DB: db}

	r := mux.NewRouter().StrictSlash(true)

	r.Handle("/signin", authentication.SignIn(env)).Methods("POST")

	r.Handle("/categories", authentication.Authenticate(category.CategoryList(env))).Methods("GET")
	r.Handle("/categories/create", category.CategoryCreate(env)).Methods("POST")
	r.Handle("/categories/{id:[0-9]+}", category.CategoryShow(env)).Methods("GET")
	r.Handle("/categories/{id:[0-9]+}", category.CategoryUpdate(env)).Methods("PUT")
	r.Handle("/categories/{id:[0-9]+}", category.CategoryDelete(env)).Methods("DELETE")

	// r.Handle("/countries", authentication.Authenticate(country.CountryList(env))).Methods("GET")
	// r.Handle("/countries/create", country.CountryCreate(env)).Methods("POST")
	// r.Handle("/countries/{id:[0-9]+}", country.CountryGet(env)).Methods("GET")
	// r.Handle("/countries/{id:[0-9]+}", country.CountryUpdate(env)).Methods("PUT")
	// r.Handle("/countries/{id:[0-9]+}", country.CountryDelete(env)).Methods("DELETE")

	return r
}
