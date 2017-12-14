package config

import (
	"github.com/jmoiron/sqlx"
)

type Env struct {
	DB *sqlx.DB
}

const DatasourceName = "user=wms password=wms dbname=wms sslmode=disable"

var JWTSigningKey = []byte("bpi-wms")
