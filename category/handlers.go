package category

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/kyawthanttin/bpi-wms/config"
	"github.com/kyawthanttin/bpi-wms/webutil"
)

func CategoryList(env *config.Env) http.Handler {
	return webutil.ListRecords(env, func(db *sql.DB, search interface{}) (interface{}, error) {
		return ListCategories(db, search.(string))
	})
}

func CategoryShow(env *config.Env) http.Handler {
	return webutil.GetRecord(env, func(db *sql.DB, id interface{}) (interface{}, error) {
		return GetCategory(db, id.(int))
	})
}

func CategoryCreate(env *config.Env) http.Handler {
	return webutil.CreateRecord(env, Category{}, func(db *sql.DB, data []byte) (interface{}, error) {
		category := Category{}
		json.Unmarshal(data, &category)
		return CreateCategory(db, category)
	})
}

func CategoryUpdate(env *config.Env) http.Handler {
	return webutil.UpdateRecord(env, Category{}, func(db *sql.DB, id interface{}, data []byte) (interface{}, error) {
		category := Category{}
		json.Unmarshal(data, &category)
		return UpdateCategory(db, id.(int), category)
	})
}

func CategoryDelete(env *config.Env) http.Handler {
	return webutil.DeleteRecord(env, func(db *sql.DB, id interface{}) error {
		return DeleteCategory(db, id.(int))
	})
}
