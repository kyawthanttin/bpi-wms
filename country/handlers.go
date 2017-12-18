package country

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/config"
	"github.com/kyawthanttin/bpi-wms/webutil"
	validator "gopkg.in/go-playground/validator.v9"
)

func CountryList(env *config.Env) http.Handler {
	return webutil.ListRecords(env, func(db *sqlx.DB, search interface{}) (interface{}, error) {
		return ListCountries(db, search.(string))
	})
}

func CountryShow(env *config.Env) http.Handler {
	return webutil.GetRecord(env, func(db *sqlx.DB, id interface{}) (interface{}, error) {
		return GetCountry(db, id.(int))
	})
}

func CountryCreate(env *config.Env) http.Handler {
	return webutil.CreateRecord(env, Country{}, func(db *sqlx.DB, validate *validator.Validate, byteData []byte) (interface{}, error) {
		data := Country{}
		json.Unmarshal(byteData, &data)
		return CreateCountry(db, validate, data)
	})
}

func CountryUpdate(env *config.Env) http.Handler {
	return webutil.UpdateRecord(env, Country{}, func(db *sqlx.DB, validate *validator.Validate, id interface{}, byteData []byte) (interface{}, error) {
		data := Country{}
		json.Unmarshal(byteData, &data)
		return UpdateCountry(db, validate, id.(int), data)
	})
}

func CountryDelete(env *config.Env) http.Handler {
	return webutil.DeleteRecord(env, func(db *sqlx.DB, id interface{}) error {
		return DeleteCountry(db, id.(int))
	})
}
