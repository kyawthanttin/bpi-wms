package unitofmeasurement

import (
	"errors"
	"strconv"

	"gopkg.in/go-playground/validator.v9"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/dbutil"
	"github.com/kyawthanttin/bpi-wms/validation"
)

type UnitOfMeasurement struct {
	Id   int    `json:"id" validate:"-"`
	Abbr string `json:"abbr" dbop:"iu" validate:"strmin=1,strmax=4,alphanum"`
	Name string `json:"name" dbop:"iu" validate:"strmin=1,strmax=30,alphanumspecial"`
}

func ListUnitOfMeasurements(db *sqlx.DB, search string) ([]UnitOfMeasurement, error) {
	results := []UnitOfMeasurement{}
	var err error

	if search != "" {
		s := UnitOfMeasurement{Name: search, Abbr: search}
		nstmt, _ := db.PrepareNamed("SELECT id, abbr, name FROM UnitOfMeasurement WHERE UPPER(name) LIKE CONCAT('%', UPPER(:name), '%') " +
			"OR UPPER(abbr) LIKE CONCAT('%', UPPER(:abbr), '%') ORDER BY name LIMIT " + strconv.Itoa(dbutil.MaxResults))
		err = nstmt.Select(&results, s)
	} else {
		err = db.Select(&results, "SELECT id, abbr, name FROM UnitOfMeasurement ORDER BY name LIMIT "+strconv.Itoa(dbutil.MaxResults))
	}
	return results, err
}

func GetUnitOfMeasurement(db *sqlx.DB, id int) (UnitOfMeasurement, error) {
	result := UnitOfMeasurement{}
	err := db.Get(&result, "SELECT id, abbr, name FROM UnitOfMeasurement WHERE id = $1", id)
	return result, err
}

func CreateUnitOfMeasurement(db *sqlx.DB, validate *validator.Validate, data UnitOfMeasurement) (UnitOfMeasurement, error) {
	if err := validate.Struct(data); err != nil {
		return UnitOfMeasurement{}, validation.DescribeErrors(err.(validator.ValidationErrors))
	}
	if exist, _ := dbutil.IsExist(db, "UnitOfMeasurement", "abbr", data.Abbr); exist {
		return UnitOfMeasurement{}, errors.New("Same unit-of-measurment already exists")
	}
	id, err := dbutil.Insert(db, "UnitOfMeasurement", &data)
	if err != nil {
		return UnitOfMeasurement{}, err
	}
	return GetUnitOfMeasurement(db, id.(int))
}

func UpdateUnitOfMeasurement(db *sqlx.DB, validate *validator.Validate, id int, data UnitOfMeasurement) (UnitOfMeasurement, error) {
	if exist, _ := dbutil.IsExist(db, "UnitOfMeasurement", "id", id); !exist {
		return UnitOfMeasurement{}, errors.New("No such unit-of-measurement")
	}
	if err := validate.Struct(data); err != nil {
		return UnitOfMeasurement{}, validation.DescribeErrors(err.(validator.ValidationErrors))
	}
	if exist, _ := dbutil.IsExistExcept(db, "UnitOfMeasurement", id, "abbr", data.Abbr); exist {
		return UnitOfMeasurement{}, errors.New("Same unit-of-measurement already exists")
	}
	err := dbutil.Update(db, "UnitOfMeasurement", &data, &UnitOfMeasurement{Id: id})
	if err != nil {
		return UnitOfMeasurement{}, err
	}
	return GetUnitOfMeasurement(db, id)
}

func DeleteUnitOfMeasurement(db *sqlx.DB, id int) error {
	if exist, _ := dbutil.IsExist(db, "UnitOfMeasurement", "id", id); !exist {
		return errors.New("No such unit-of-measurement")
	}
	_, err := db.Exec("DELETE FROM UnitOfMeasurement WHERE id = $1", id)
	return err
}
