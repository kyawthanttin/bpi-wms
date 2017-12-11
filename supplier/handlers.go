package supplier

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/kyawthanttin/bpi-wms/config"
	"github.com/kyawthanttin/bpi-wms/webutil"
)

func SupplierList(env *config.Env) http.Handler {
	return webutil.ListRecords(env, func(db *sql.DB, search interface{}) (interface{}, error) {
		return ListSuppliers(db, search.(string))
	})
}

func SupplierShow(env *config.Env) http.Handler {
	return webutil.GetRecord(env, func(db *sql.DB, id interface{}) (interface{}, error) {
		return GetSupplier(db, id.(int))
	})
}

func SupplierCreate(env *config.Env) http.Handler {
	return webutil.CreateRecord(env, Supplier{}, func(db *sql.DB, byteData []byte) (interface{}, error) {
		data := Supplier{}
		json.Unmarshal(byteData, &data)
		return CreateSupplier(db, data)
	})
}

func SupplierUpdate(env *config.Env) http.Handler {
	return webutil.UpdateRecord(env, Supplier{}, func(db *sql.DB, id interface{}, byteData []byte) (interface{}, error) {
		data := Supplier{}
		json.Unmarshal(byteData, &data)
		return UpdateSupplier(db, id.(int), data)
	})
}

func SupplierDelete(env *config.Env) http.Handler {
	return webutil.DeleteRecord(env, func(db *sql.DB, id interface{}) error {
		return DeleteSupplier(db, id.(int))
	})
}
