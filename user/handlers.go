package user

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

func UserList(env *config.Env) http.Handler {
	return webutil.ListRecords(env, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, search interface{}) (interface{}, error) {
		if err := authentication.CheckPrivilegeByRoles(r, []string{authentication.AdminRolename}); err != nil {
			return nil, validation.NewErrorWithHttpStatus(http.StatusUnauthorized, err.Error())
		}
		return ListUsers(db, search.(string))
	})
}

func UserShow(env *config.Env) http.Handler {
	return webutil.GetRecord(env, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, id interface{}) (interface{}, error) {
		if err := authentication.CheckPrivilegeByRoles(r, []string{authentication.AdminRolename}); err != nil {
			if err = authentication.CheckPrivilegeById(r, id.(int)); err != nil {
				return nil, validation.NewErrorWithHttpStatus(http.StatusUnauthorized, err.Error())
			}
		}
		return GetUser(db, id.(int))
	})
}

func UserCreate(env *config.Env) http.Handler {
	return webutil.CreateRecord(env, User{}, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, validate *validator.Validate, byteData []byte) (interface{}, error) {
		if err := authentication.CheckPrivilegeByRoles(r, []string{authentication.AdminRolename}); err != nil {
			return nil, validation.NewErrorWithHttpStatus(http.StatusUnauthorized, err.Error())
		}
		data := User{}
		json.Unmarshal(byteData, &data)
		return CreateUser(db, validate, data)
	})
}

func PasswordChange(env *config.Env) http.Handler {
	return webutil.UpdateRecord(env, User{}, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, validate *validator.Validate, id interface{}, byteData []byte) (interface{}, error) {
		if err := authentication.CheckPrivilegeByRoles(r, []string{authentication.AdminRolename}); err != nil {
			if err = authentication.CheckPrivilegeById(r, id.(int)); err != nil {
				return nil, validation.NewErrorWithHttpStatus(http.StatusUnauthorized, err.Error())
			}
		}
		data := User{}
		json.Unmarshal(byteData, &data)
		return ChangePassword(db, validate, id.(int), data.Password)
	})
}

func UserUpdate(env *config.Env) http.Handler {
	return webutil.UpdateRecord(env, User{}, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, validate *validator.Validate, id interface{}, byteData []byte) (interface{}, error) {
		if err := authentication.CheckPrivilegeByRoles(r, []string{authentication.AdminRolename}); err != nil {
			if err = authentication.CheckPrivilegeById(r, id.(int)); err != nil {
				return nil, validation.NewErrorWithHttpStatus(http.StatusUnauthorized, err.Error())
			}
		}
		data := User{}
		json.Unmarshal(byteData, &data)
		return UpdateUser(db, validate, id.(int), data)
	})
}

func UserDelete(env *config.Env) http.Handler {
	return webutil.DeleteRecord(env, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, id interface{}) error {
		if err := authentication.CheckPrivilegeByRoles(r, []string{authentication.AdminRolename}); err != nil {
			return validation.NewErrorWithHttpStatus(http.StatusUnauthorized, err.Error())
		}
		return DeleteUser(db, id.(int))
	})
}
