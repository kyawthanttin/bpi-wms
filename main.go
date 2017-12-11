package main

import (
	"bpi-wms/models"
	"log"
	"net/http"
)

type Env struct {
	db models.Datastore
}

func main() {
	db, err := models.NewDB("postgres://wms:wms@localhost:5432/bpi_wms")
	if err != nil {
		log.Panic(err)
	}

	env := &Env{db}

	// http.HandleFunc("/categories", env.)
}

func (env *Env) categoriesIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		RespondWithError(w, http.StatusText(405), 405)
		return
	}
	results, err := env.db.AllCategories()
	if err != nil {
		RespondWithError(w, http.StatusText(500), err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, results)
}
