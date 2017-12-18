package supplier

import (
	"errors"
	"strconv"

	"gopkg.in/go-playground/validator.v9"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/dbutil"
	"github.com/kyawthanttin/bpi-wms/validation"
)

type Supplier struct {
	Id      int    `json:"id" validate:"-"`
	Name    string `json:"name" dbop:"iu" validate:"strmin=1,strmax=50,alphanumspecial"`
	Address string `json:"address" dbop:"iu" validate:"strmax=250"`
}

func ListSuppliers(db *sqlx.DB, search string) ([]Supplier, error) {
	results := []Supplier{}
	var err error

	if search != "" {
		s := Supplier{Name: search}
		nstmt, _ := db.PrepareNamed("SELECT id, name, address FROM Supplier WHERE UPPER(name) LIKE CONCAT('%', UPPER(:name), '%') ORDER BY name LIMIT " + strconv.Itoa(dbutil.MaxResults))
		err = nstmt.Select(&results, s)
	} else {
		err = db.Select(&results, "SELECT id, name, address FROM Supplier ORDER BY name LIMIT "+strconv.Itoa(dbutil.MaxResults))
	}
	return results, err
}

func GetSupplier(db *sqlx.DB, id int) (Supplier, error) {
	result := Supplier{}
	err := db.Get(&result, "SELECT id, name, address FROM Supplier WHERE id = $1", id)
	return result, err
}

func CreateSupplier(db *sqlx.DB, validate *validator.Validate, data Supplier) (Supplier, error) {
	if err := validate.Struct(data); err != nil {
		return Supplier{}, validation.DescribeErrors(err.(validator.ValidationErrors))
	}
	if exist, _ := dbutil.IsExist(db, "Supplier", "name", data.Name); exist {
		return Supplier{}, errors.New("Same supplier already exists")
	}
	id, err := dbutil.Insert(db, "Supplier", &data)
	if err != nil {
		return Supplier{}, err
	}
	return GetSupplier(db, id.(int))
}

func UpdateSupplier(db *sqlx.DB, validate *validator.Validate, id int, data Supplier) (Supplier, error) {
	if exist, _ := dbutil.IsExist(db, "Supplier", "id", id); !exist {
		return Supplier{}, errors.New("No such supplier")
	}
	if err := validate.Struct(data); err != nil {
		return Supplier{}, validation.DescribeErrors(err.(validator.ValidationErrors))
	}
	if exist, _ := dbutil.IsExistExcept(db, "Supplier", id, "name", data.Name); exist {
		return Supplier{}, errors.New("Same supplier already exists")
	}
	err := dbutil.Update(db, "Supplier", &data, &Supplier{Id: id})
	if err != nil {
		return Supplier{}, err
	}
	return GetSupplier(db, id)
}

func DeleteSupplier(db *sqlx.DB, id int) error {
	if exist, _ := dbutil.IsExist(db, "Supplier", "id", id); !exist {
		return errors.New("No such supplier")
	}
	_, err := db.Exec("DELETE FROM Supplier WHERE id = $1", id)
	return err
}
