package country

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/kyawthanttin/bpi-wms/config"
	"github.com/kyawthanttin/bpi-wms/webutil"
)

func CountryList(env *config.Env) http.Handler {
	return webutil.ListRecords(env, func(db *sql.DB, search interface{}) (interface{}, error) {
		return ListCountries(db, search.(string))
	})
}

func CountryShow(env *config.Env) http.Handler {
	return webutil.GetRecord(env, func(db *sql.DB, id interface{}) (interface{}, error) {
		return GetCountry(db, id.(int))
	})
}

func CountryCreate(env *config.Env) http.Handler {
	return webutil.CreateRecord(env, Country{}, func(db *sql.DB, byteData []byte) (interface{}, error) {
		data := Country{}
		json.Unmarshal(byteData, &data)
		return CreateCountry(db, data)
	})
}

func CountryUpdate(env *config.Env) http.Handler {
	return webutil.UpdateRecord(env, Country{}, func(db *sql.DB, id interface{}, byteData []byte) (interface{}, error) {
		data := Country{}
		json.Unmarshal(byteData, &data)
		return UpdateCountry(db, id.(int), data)
	})
}

func CountryDelete(env *config.Env) http.Handler {
	return webutil.DeleteRecord(env, func(db *sql.DB, id interface{}) error {
		return DeleteCountry(db, id.(int))
	})
}
