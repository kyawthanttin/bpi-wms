package supplier

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/config"
	"github.com/kyawthanttin/bpi-wms/webutil"
)

func SupplierList(env *config.Env) http.Handler {
	return webutil.ListRecords(env, func(db *sqlx.DB, search interface{}) (interface{}, error) {
		return ListSuppliers(db, search.(string))
	})
}

func SupplierShow(env *config.Env) http.Handler {
	return webutil.GetRecord(env, func(db *sqlx.DB, id interface{}) (interface{}, error) {
		return GetSupplier(db, id.(int))
	})
}

func SupplierCreate(env *config.Env) http.Handler {
	return webutil.CreateRecord(env, Supplier{}, func(db *sqlx.DB, byteData []byte) (interface{}, error) {
		data := Supplier{}
		json.Unmarshal(byteData, &data)
		return CreateSupplier(db, data)
	})
}

func SupplierUpdate(env *config.Env) http.Handler {
	return webutil.UpdateRecord(env, Supplier{}, func(db *sqlx.DB, id interface{}, byteData []byte) (interface{}, error) {
		data := Supplier{}
		json.Unmarshal(byteData, &data)
		return UpdateSupplier(db, id.(int), data)
	})
}

func SupplierDelete(env *config.Env) http.Handler {
	return webutil.DeleteRecord(env, func(db *sqlx.DB, id interface{}) error {
		return DeleteSupplier(db, id.(int))
	})
}
