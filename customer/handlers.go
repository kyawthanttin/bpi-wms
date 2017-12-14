package customer

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/config"
	"github.com/kyawthanttin/bpi-wms/webutil"
)

func CustomerList(env *config.Env) http.Handler {
	return webutil.ListRecords(env, func(db *sqlx.DB, search interface{}) (interface{}, error) {
		return ListCustomers(db, search.(string))
	})
}

func CustomerShow(env *config.Env) http.Handler {
	return webutil.GetRecord(env, func(db *sqlx.DB, id interface{}) (interface{}, error) {
		return GetCustomer(db, id.(int))
	})
}

func CustomerCreate(env *config.Env) http.Handler {
	return webutil.CreateRecord(env, Customer{}, func(db *sqlx.DB, byteData []byte) (interface{}, error) {
		data := Customer{}
		json.Unmarshal(byteData, &data)
		return CreateCustomer(db, data)
	})
}

func CustomerUpdate(env *config.Env) http.Handler {
	return webutil.UpdateRecord(env, Customer{}, func(db *sqlx.DB, id interface{}, byteData []byte) (interface{}, error) {
		data := Customer{}
		json.Unmarshal(byteData, &data)
		return UpdateCustomer(db, id.(int), data)
	})
}

func CustomerDelete(env *config.Env) http.Handler {
	return webutil.DeleteRecord(env, func(db *sqlx.DB, id interface{}) error {
		return DeleteCustomer(db, id.(int))
	})
}
