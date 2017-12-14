package user

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/config"
	"github.com/kyawthanttin/bpi-wms/webutil"
)

func UserList(env *config.Env) http.Handler {
	return webutil.ListRecords(env, func(db *sqlx.DB, search interface{}) (interface{}, error) {
		return ListUsers(db, search.(string))
	})
}

func UserShow(env *config.Env) http.Handler {
	return webutil.GetRecord(env, func(db *sqlx.DB, id interface{}) (interface{}, error) {
		return GetUser(db, id.(int))
	})
}

func UserCreate(env *config.Env) http.Handler {
	return webutil.CreateRecord(env, User{}, func(db *sqlx.DB, byteData []byte) (interface{}, error) {
		data := User{}
		json.Unmarshal(byteData, &data)
		return CreateUser(db, data)
	})
}

func PasswordChange(env *config.Env) http.Handler {
	return webutil.UpdateRecord(env, User{}, func(db *sqlx.DB, id interface{}, byteData []byte) (interface{}, error) {
		data := User{}
		json.Unmarshal(byteData, &data)
		return ChangePassword(db, id.(int), data.Password)
	})
}

func UserUpdate(env *config.Env) http.Handler {
	return webutil.UpdateRecord(env, User{}, func(db *sqlx.DB, id interface{}, byteData []byte) (interface{}, error) {
		data := User{}
		json.Unmarshal(byteData, &data)
		return UpdateUser(db, id.(int), data)
	})
}

func UserDelete(env *config.Env) http.Handler {
	return webutil.DeleteRecord(env, func(db *sqlx.DB, id interface{}) error {
		return DeleteUser(db, id.(int))
	})
}
