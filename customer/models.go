package customer

import (
	"errors"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/kyawthanttin/bpi-wms/dbutil"
	"github.com/kyawthanttin/bpi-wms/validation"
	validator "gopkg.in/go-playground/validator.v9"
)

type Customer struct {
	Id      int    `json:"id" validate:"-"`
	Name    string `json:"name" dbop:"iu" validate:"strmin=1,strmax=50,alphanumspecial"`
	Address string `json:"address" dbop:"iu" validate:"strmax=250"`
}

func ListCustomers(db *sqlx.DB, search string) ([]Customer, error) {
	results := []Customer{}
	var err error

	if search != "" {
		s := Customer{Name: search}
		nstmt, _ := db.PrepareNamed("SELECT id, name, address FROM Customer WHERE UPPER(name) LIKE CONCAT('%', UPPER(:name), '%') ORDER BY name LIMIT " + strconv.Itoa(dbutil.MaxResults))
		err = nstmt.Select(&results, s)
	} else {
		err = db.Select(&results, "SELECT id, name, address FROM Customer ORDER BY name LIMIT "+strconv.Itoa(dbutil.MaxResults))
	}
	return results, err
}

func GetCustomer(db *sqlx.DB, id int) (Customer, error) {
	result := Customer{}
	err := db.Get(&result, "SELECT id, name, address FROM Customer WHERE id = $1", id)
	return result, err
}

func CreateCustomer(db *sqlx.DB, validate *validator.Validate, data Customer) (Customer, error) {
	if err := validate.Struct(data); err != nil {
		return Customer{}, validation.DescribeErrors(err.(validator.ValidationErrors))
	}
	if exist, _ := dbutil.IsExist(db, "Customer", "name", data.Name); exist {
		return Customer{}, errors.New("Same customer already exists")
	}
	id, err := dbutil.Insert(db, "Customer", &data)
	if err != nil {
		return Customer{}, err
	}
	return GetCustomer(db, id.(int))
}

func UpdateCustomer(db *sqlx.DB, validate *validator.Validate, id int, data Customer) (Customer, error) {
	if exist, _ := dbutil.IsExist(db, "Customer", "id", id); !exist {
		return Customer{}, errors.New("No such customer")
	}
	if err := validate.Struct(data); err != nil {
		return Customer{}, validation.DescribeErrors(err.(validator.ValidationErrors))
	}
	if exist, _ := dbutil.IsExistExcept(db, "Customer", id, "name", data.Name); exist {
		return Customer{}, errors.New("Same customer already exists")
	}
	err := dbutil.Update(db, "Customer", &data, &Customer{Id: id})
	if err != nil {
		return Customer{}, err
	}
	return GetCustomer(db, id)
}

func DeleteCustomer(db *sqlx.DB, id int) error {
	if exist, _ := dbutil.IsExist(db, "Customer", "id", id); !exist {
		return errors.New("No such customer")
	}
	_, err := db.Exec("DELETE FROM Customer WHERE id = $1", id)
	return err
}
