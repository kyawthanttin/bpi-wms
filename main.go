package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func main() {

	router := NewRouter()

	// these two lines are important in order to allow access from the front-end side to the methods
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})

	// launch server with CORS validations
	log.Fatal(http.ListenAndServe(":3000",
		handlers.CORS(allowedOrigins, allowedMethods)(router)))
}
