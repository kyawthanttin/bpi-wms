package unitofmeasurement

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

func UnitOfMeasurementList(env *config.Env) http.Handler {
	return webutil.ListRecords(env, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, search interface{}) (interface{}, error) {
		return ListUnitOfMeasurements(db, search.(string))
	})
}

func UnitOfMeasurementShow(env *config.Env) http.Handler {
	return webutil.GetRecord(env, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, id interface{}) (interface{}, error) {
		if err := authentication.CheckPrivilegeByRoles(r, []string{authentication.AdminRolename}); err != nil {
			return nil, validation.NewErrorWithHttpStatus(http.StatusUnauthorized, err.Error())
		}
		return GetUnitOfMeasurement(db, id.(int))
	})
}

func UnitOfMeasurementCreate(env *config.Env) http.Handler {
	return webutil.CreateRecord(env, UnitOfMeasurement{}, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, validate *validator.Validate, byteData []byte) (interface{}, error) {
		if err := authentication.CheckPrivilegeByRoles(r, []string{authentication.AdminRolename}); err != nil {
			return nil, validation.NewErrorWithHttpStatus(http.StatusUnauthorized, err.Error())
		}
		data := UnitOfMeasurement{}
		json.Unmarshal(byteData, &data)
		return CreateUnitOfMeasurement(db, validate, data)
	})
}

func UnitOfMeasurementUpdate(env *config.Env) http.Handler {
	return webutil.UpdateRecord(env, UnitOfMeasurement{}, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, validate *validator.Validate, id interface{}, byteData []byte) (interface{}, error) {
		if err := authentication.CheckPrivilegeByRoles(r, []string{authentication.AdminRolename}); err != nil {
			return nil, validation.NewErrorWithHttpStatus(http.StatusUnauthorized, err.Error())
		}
		data := UnitOfMeasurement{}
		json.Unmarshal(byteData, &data)
		return UpdateUnitOfMeasurement(db, validate, id.(int), data)
	})
}

func UnitOfMeasurementDelete(env *config.Env) http.Handler {
	return webutil.DeleteRecord(env, func(w http.ResponseWriter, r *http.Request, db *sqlx.DB, id interface{}) error {
		if err := authentication.CheckPrivilegeByRoles(r, []string{authentication.AdminRolename}); err != nil {
			return validation.NewErrorWithHttpStatus(http.StatusUnauthorized, err.Error())
		}
		return DeleteUnitOfMeasurement(db, id.(int))
	})
}
