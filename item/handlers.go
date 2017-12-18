package item

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/config"
	"github.com/kyawthanttin/bpi-wms/webutil"
	validator "gopkg.in/go-playground/validator.v9"
)

func ItemList(env *config.Env) http.Handler {
	return webutil.ListRecords(env, func(db *sqlx.DB, search interface{}) (interface{}, error) {
		return ListItems(db, search.(string))
	})
}

func ItemShow(env *config.Env) http.Handler {
	return webutil.GetRecord(env, func(db *sqlx.DB, id interface{}) (interface{}, error) {
		return GetItem(db, id.(int))
	})
}

func ItemCreate(env *config.Env) http.Handler {
	return webutil.CreateRecord(env, Item{}, func(db *sqlx.DB, validate *validator.Validate, byteData []byte) (interface{}, error) {
		data := Item{}
		json.Unmarshal(byteData, &data)
		return CreateItem(db, validate, data)
	})
}

func ItemUpdate(env *config.Env) http.Handler {
	return webutil.UpdateRecord(env, Item{}, func(db *sqlx.DB, validate *validator.Validate, id interface{}, byteData []byte) (interface{}, error) {
		data := Item{}
		json.Unmarshal(byteData, &data)
		return UpdateItem(db, validate, id.(int), data)
	})
}

func ItemDelete(env *config.Env) http.Handler {
	return webutil.DeleteRecord(env, func(db *sqlx.DB, id interface{}) error {
		return DeleteItem(db, id.(int))
	})
}
