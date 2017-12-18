package config

import (
	"github.com/jmoiron/sqlx"
	"gopkg.in/go-playground/validator.v9"
)

type Env struct {
	DB       *sqlx.DB
	Validate *validator.Validate
}

const DatasourceName = "user=wms password=wms dbname=wms sslmode=disable"

var JWTSigningKey = []byte("bpi-wms")
