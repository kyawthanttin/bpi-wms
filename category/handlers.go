package category

import (
	"encoding/json"
	"net/http"

	"gopkg.in/go-playground/validator.v9"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/config"
	"github.com/kyawthanttin/bpi-wms/webutil"
)

func CategoryList(env *config.Env) http.Handler {
	return webutil.ListRecords(env, func(db *sqlx.DB, search interface{}) (interface{}, error) {
		return ListCategories(db, search.(string))
	})
}

func CategoryShow(env *config.Env) http.Handler {
	return webutil.GetRecord(env, func(db *sqlx.DB, id interface{}) (interface{}, error) {
		return GetCategory(db, id.(int))
	})
}

func CategoryCreate(env *config.Env) http.Handler {
	return webutil.CreateRecord(env, Category{}, func(db *sqlx.DB, validate *validator.Validate, byteData []byte) (interface{}, error) {
		data := Category{}
		json.Unmarshal(byteData, &data)
		return CreateCategory(db, validate, data)
	})
}

func CategoryUpdate(env *config.Env) http.Handler {
	return webutil.UpdateRecord(env, Category{}, func(db *sqlx.DB, validate *validator.Validate, id interface{}, byteData []byte) (interface{}, error) {
		data := Category{}
		json.Unmarshal(byteData, &data)
		return UpdateCategory(db, validate, id.(int), data)
	})
}

func CategoryDelete(env *config.Env) http.Handler {
	return webutil.DeleteRecord(env, func(db *sqlx.DB, id interface{}) error {
		return DeleteCategory(db, id.(int))
	})
}
