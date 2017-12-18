package item

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/authentication"
	"github.com/kyawthanttin/bpi-wms/config"
	"github.com/kyawthanttin/bpi-wms/validation"
	"github.com/kyawthanttin/bpi-wms/webutil"
	validator "gopkg.in/go-playground/validator.v9"
)

func ItemList(env *config.Env) http.Handler {
	return webutil.ListRecords(env, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, search interface{}) (interface{}, error) {
		return ListItems(db, search.(string))
	})
}

func ItemShow(env *config.Env) http.Handler {
	return webutil.GetRecord(env, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, id interface{}) (interface{}, error) {
		if err := authentication.CheckPrivilegeByRoles(r, []string{authentication.AdminRolename}); err != nil {
			return nil, validation.NewErrorWithHttpStatus(http.StatusUnauthorized, err.Error())
		}
		return GetItem(db, id.(int))
	})
}

func ItemCreate(env *config.Env) http.Handler {
	return webutil.CreateRecord(env, Item{}, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, validate *validator.Validate, byteData []byte) (interface{}, error) {
		if err := authentication.CheckPrivilegeByRoles(r, []string{authentication.AdminRolename}); err != nil {
			return nil, validation.NewErrorWithHttpStatus(http.StatusUnauthorized, err.Error())
		}
		data := Item{}
		json.Unmarshal(byteData, &data)
		return CreateItem(db, validate, data)
	})
}

func ItemUpdate(env *config.Env) http.Handler {
	return webutil.UpdateRecord(env, Item{}, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, validate *validator.Validate, id interface{}, byteData []byte) (interface{}, error) {
		if err := authentication.CheckPrivilegeByRoles(r, []string{authentication.AdminRolename}); err != nil {
			return nil, validation.NewErrorWithHttpStatus(http.StatusUnauthorized, err.Error())
		}
		data := Item{}
		json.Unmarshal(byteData, &data)
		return UpdateItem(db, validate, id.(int), data)
	})
}

func ItemDelete(env *config.Env) http.Handler {
	return webutil.DeleteRecord(env, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, id interface{}) error {
		if err := authentication.CheckPrivilegeByRoles(r, []string{authentication.AdminRolename}); err != nil {
			return validation.NewErrorWithHttpStatus(http.StatusUnauthorized, err.Error())
		}
		return DeleteItem(db, id.(int))
	})
}
