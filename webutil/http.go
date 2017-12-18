package webutil

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"gopkg.in/go-playground/validator.v9"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/config"
)

func RespondWithJSON(w http.ResponseWriter, code int, data interface{}) {
	response, _ := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func RespondWithErrorType(w http.ResponseWriter, err error) {
	switch err {
	case sql.ErrNoRows:
		RespondWithError(w, http.StatusNotFound, "There is no such data")
	default:
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
}

func ListRecords(env *config.Env, listFunc func(*sqlx.DB, interface{}) (interface{}, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
			return
		}
		results, err := listFunc(env.DB, r.FormValue("search"))
		if err != nil {
			RespondWithErrorType(w, err)
			return
		}
		RespondWithJSON(w, http.StatusOK, results)
	})
}

func GetRecord(env *config.Env, getFunc func(*sqlx.DB, interface{}) (interface{}, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
			return
		}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid Id")
			return
		}

		result, err := getFunc(env.DB, id)
		if err != nil {
			RespondWithErrorType(w, err)
			return
		}
		RespondWithJSON(w, http.StatusOK, result)
	})
}

func CreateRecord(env *config.Env, data interface{}, createFunc func(*sqlx.DB, *validator.Validate, []byte) (interface{}, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
			return
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&data); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		byteData, _ := json.Marshal(data)
		created, err := createFunc(env.DB, env.Validate, byteData)
		if err != nil {
			RespondWithErrorType(w, err)
			return
		}
		RespondWithJSON(w, http.StatusCreated, created)
	})
}

func UpdateRecord(env *config.Env, data interface{}, updateFunc func(*sqlx.DB, *validator.Validate, interface{}, []byte) (interface{}, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			RespondWithError(w, http.StatusBadRequest, "Method Not Allowed")
			return
		}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid Id")
			return
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&data); err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()
		byteData, _ := json.Marshal(data)
		updated, err := updateFunc(env.DB, env.Validate, id, byteData)
		if err != nil {
			RespondWithErrorType(w, err)
			return
		}
		RespondWithJSON(w, http.StatusOK, updated)
	})
}

func DeleteRecord(env *config.Env, deleteFunc func(*sqlx.DB, interface{}) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			RespondWithError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
			return
		}
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Invalid Id")
			return
		}
		if err := deleteFunc(env.DB, id); err != nil {
			RespondWithErrorType(w, err)
			return
		}
		RespondWithJSON(w, http.StatusOK, "Record Deleted")
	})
}
