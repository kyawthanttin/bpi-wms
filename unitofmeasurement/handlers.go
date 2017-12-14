package unitofmeasurement

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/config"
	"github.com/kyawthanttin/bpi-wms/webutil"
)

func UnitOfMeasurementList(env *config.Env) http.Handler {
	return webutil.ListRecords(env, func(db *sqlx.DB, search interface{}) (interface{}, error) {
		return ListUnitOfMeasurements(db, search.(string))
	})
}

func UnitOfMeasurementShow(env *config.Env) http.Handler {
	return webutil.GetRecord(env, func(db *sqlx.DB, id interface{}) (interface{}, error) {
		return GetUnitOfMeasurement(db, id.(int))
	})
}

func UnitOfMeasurementCreate(env *config.Env) http.Handler {
	return webutil.CreateRecord(env, UnitOfMeasurement{}, func(db *sqlx.DB, byteData []byte) (interface{}, error) {
		data := UnitOfMeasurement{}
		json.Unmarshal(byteData, &data)
		return CreateUnitOfMeasurement(db, data)
	})
}

func UnitOfMeasurementUpdate(env *config.Env) http.Handler {
	return webutil.UpdateRecord(env, UnitOfMeasurement{}, func(db *sqlx.DB, id interface{}, byteData []byte) (interface{}, error) {
		data := UnitOfMeasurement{}
		json.Unmarshal(byteData, &data)
		return UpdateUnitOfMeasurement(db, id.(int), data)
	})
}

func UnitOfMeasurementDelete(env *config.Env) http.Handler {
	return webutil.DeleteRecord(env, func(db *sqlx.DB, id interface{}) error {
		return DeleteUnitOfMeasurement(db, id.(int))
	})
}
