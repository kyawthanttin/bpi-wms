package customer

import (
	"encoding/json"
	"net/http"

	"gopkg.in/go-playground/validator.v9"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/authentication"
	"github.com/kyawthanttin/bpi-wms/config"
	"github.com/kyawthanttin/bpi-wms/validation"
	"github.com/kyawthanttin/bpi-wms/webutil"
)

func CustomerList(env *config.Env) http.Handler {
	return webutil.ListRecords(env, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, search interface{}) (interface{}, error) {
		return ListCustomers(db, search.(string))
	})
}

func CustomerShow(env *config.Env) http.Handler {
	return webutil.GetRecord(env, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, id interface{}) (interface{}, error) {
		if err := authentication.CheckPrivilegeByRoles(r, []string{authentication.AdminRolename}); err != nil {
			return nil, validation.NewErrorWithHttpStatus(http.StatusUnauthorized, err.Error())
		}
		return GetCustomer(db, id.(int))
	})
}

func CustomerCreate(env *config.Env) http.Handler {
	return webutil.CreateRecord(env, Customer{}, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, validate *validator.Validate, byteData []byte) (interface{}, error) {
		if err := authentication.CheckPrivilegeByRoles(r, []string{authentication.AdminRolename}); err != nil {
			return nil, validation.NewErrorWithHttpStatus(http.StatusUnauthorized, err.Error())
		}
		data := Customer{}
		json.Unmarshal(byteData, &data)
		return CreateCustomer(db, validate, data)
	})
}

func CustomerUpdate(env *config.Env) http.Handler {
	return webutil.UpdateRecord(env, Customer{}, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, validate *validator.Validate, id interface{}, byteData []byte) (interface{}, error) {
		if err := authentication.CheckPrivilegeByRoles(r, []string{authentication.AdminRolename}); err != nil {
			return nil, validation.NewErrorWithHttpStatus(http.StatusUnauthorized, err.Error())
		}
		data := Customer{}
		json.Unmarshal(byteData, &data)
		return UpdateCustomer(db, validate, id.(int), data)
	})
}

func CustomerDelete(env *config.Env) http.Handler {
	return webutil.DeleteRecord(env, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, id interface{}) error {
		if err := authentication.CheckPrivilegeByRoles(r, []string{authentication.AdminRolename}); err != nil {
			return validation.NewErrorWithHttpStatus(http.StatusUnauthorized, err.Error())
		}
		return DeleteCustomer(db, id.(int))
	})
}
